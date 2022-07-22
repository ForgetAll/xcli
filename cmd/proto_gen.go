package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type protoFileInfo struct {
	fullPath string
	fileName string
}

var (
	protoPath  string
	outputPath = "."
	verbose    bool

	moduleName string
)

var protoGenCmd = &cobra.Command{
	Use:   "proto-gen",
	Short: "generate code by proto file",
	Long:  `generate code by proto file, just generate default empty implement code`,
	Run: func(cmd *cobra.Command, args []string) {
		if protoPath == "" {
			cmd.Println("proto file path is empty")
			return
		}

		if !checkProtoPathDirExists(protoPath) {
			cmd.Println("proto file path is not exists")
			return
		}

		if !checkProtoPathDirExists(outputPath) {
			cmd.Println("output file path is not exists")
			return
		}

		var err error
		moduleName, err = getModuleNameByGoMod()
		if err != nil {
			cmd.Println(err)
			return
		}

		if moduleName == "" {
			cmd.Println("module name is empty")
			return
		}

		if verbose {
			cmd.Println("module name: " + moduleName)
		}

		protoFiles, err := findAllProtoFile(protoPath)
		if err != nil {
			cmd.Println(err)
			return
		}

		for i := range protoFiles {
			if verbose {
				cmd.Printf("%v generate start!\n", protoFiles[i].fileName)
			}
			code, err := generateCodeByProto(protoFiles[i].fullPath)
			if err != nil {
				cmd.Println(err)
				return
			}
			err = writeCodeToFile(code, protoFiles[i].fileName)
			if err != nil {
				cmd.Printf("write code to file error: %v\n", err)
				return
			}

			if verbose {
				cmd.Printf("%v generate finish!\n", protoFiles[i].fileName)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(protoGenCmd)
	protoGenCmd.PersistentFlags().StringVarP(&protoPath, "path", "p", "", "set proto file path")
	protoGenCmd.PersistentFlags().StringVarP(&outputPath, "outputpath", "o", "", "set proto file path")
	protoGenCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "set verbose mode")
}

// checkProtoPathDirExists checks if the proto path exists.
func checkProtoPathDirExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// getModuleNameByGoMod returns the module name by go.mod file.
func getModuleNameByGoMod() (string, error) {
	if !strings.HasSuffix(protoPath, `/`) {
		protoPath += `/`
	}

	dir, err := os.ReadDir(protoPath)
	if err != nil {
		return "", err
	}

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}

		if fi.Name() == "go.mod" {
			f, err := os.Open(protoPath + fi.Name())
			if err != nil {
				return "", err
			}

			nr := bufio.NewReader(f)
			var currentLine string
			for {
				data, isPrefix, err := nr.ReadLine()
				if err == io.EOF {
					return "", nil
				}
				if err != nil {
					return "", err
				}

				if isPrefix {
					currentLine += string(data)
					continue
				} else {
					currentLine += string(data)
				}

				if currentLine == "" {
					continue
				}

				if strings.HasPrefix(currentLine, "module") {
					return strings.TrimSpace(strings.Split(currentLine, " ")[1]), nil
				}

				if !isPrefix {
					currentLine = ""
				}
			}
		}
	}

	return "", errors.New("go.mod file is not exists")
}

// findAllProtoFile find all file which name end with .proto.
func findAllProtoFile(dirPath string) ([]*protoFileInfo, error) {
	if !strings.HasSuffix(dirPath, `/`) {
		dirPath += `/`
	}

	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var protoFiles []*protoFileInfo
	for _, fi := range dir {
		if fi.IsDir() {
			files, err := findAllProtoFile(dirPath + fi.Name())
			if err != nil {
				return protoFiles, err
			}

			protoFiles = append(protoFiles, files...)
			continue
		}

		if strings.HasSuffix(fi.Name(), ".proto") {
			file := &protoFileInfo{
				fullPath: dirPath + fi.Name(),
				fileName: fi.Name(),
			}
			protoFiles = append(protoFiles, file)
		}
	}

	return protoFiles, nil
}

func writeCodeToFile(code string, fileName string) error {
	if !strings.HasSuffix(outputPath, `/`) {
		outputPath += `/`
	}

	return ioutil.WriteFile(fmt.Sprintf(outputPath+"%v_generte.go", strings.ReplaceAll(fileName, ".proto", "")), []byte(code), 0o777)
}

// generateCodeByProto generates code by proto file.
func generateCodeByProto(protoFile string) (string, error) {
	f, err := os.Open(protoFile)
	if err != nil {
		return "", err
	}

	nr := bufio.NewReader(f)
	currentLine := ""
	startRecord := false

	code := ""
	code += "package core\n\n"
	code += "import(\n"
	code += "	\"context\"\n\n"
	code += "	\"" + moduleName + "\"\n"
	code += ")\n\n"

	for {
		line, isPrefix, err := nr.ReadLine()
		if err == io.EOF {
			return code, nil
		}

		if isPrefix {
			currentLine += string(line)
			continue
		} else {
			currentLine += string(line)
		}

		if currentLine == "" {
			continue
		}

		if !startRecord && strings.Contains(currentLine, "service") {
			startRecord = true
			currentLine = ""
			continue
		}

		if !startRecord {
			continue
		}

		if strings.Contains(currentLine, "//") {
			newLine := ""
			trim := true
			for _, ch := range currentLine {
				if trim && ch == ' ' {
					continue
				}

				newLine += string(ch)
				trim = false
			}

			code += newLine + "\n"
			currentLine = ""
			continue
		}

		method, req, resp := getMethodReqRespName(currentLine)
		if method == "" || req == "" || resp == "" {
			continue
		}
		generateFuncCode(&code, method, req, resp)
		currentLine = ""
	}
}

func generateFuncCode(code *string, method string, req string, resp string) {
	*code = *code + fmt.Sprintf("func (c *Core) %v(ctx context.Context, in *rpc.%v) (*rpc.%v, error) {\n", method, req, resp)
	*code = *code + fmt.Sprintf("	resp := &rpc.%v{\n", resp)
	*code = *code + "		Response: &rpc.BaseResponse{\n"
	*code = *code + "			Result: http.StatusInternalServerError,\n"
	*code = *code + "		},\n"
	*code = *code + "	}\n\n"
	*code = *code + "	return resp, nil\n"
	*code = *code + "}\n\n"
}

// nolint
// getMethodReqRespName returns the method request and response name.
func getMethodReqRespName(line string) (method string, req string, resp string) {
	if line == "" {
		return
	}

	matchRpc := 0b000
	paramRead := 0b000
	for _, ch := range line {
		if matchRpc == 0b111 {
			switch paramRead {
			case 0b000:
				if ' ' == ch {
					if len(method) > 0 {
						paramRead = 0b001
					}
					continue
				}

				method += string(ch)
			case 0b001:
				if '(' == ch {
					if len(req) == 0 {
						paramRead = 0b010
					} else {
						paramRead = 0b011
					}
				}
			case 0b010:
				if ')' == ch {
					if len(req) > 0 {
						paramRead = 0b001
					}
					continue
				}

				req += string(ch)
			case 0b11:
				if ')' == ch {
					matchRpc = 0b000
					break
				}

				resp += string(ch)
			}

			continue
		}

		if ch == 'r' {
			matchRpc = 0b001
			continue
		}

		if ch == 'p' {
			if matchRpc == 0b001 {
				matchRpc = 0b011
				continue
			} else {
				matchRpc = 0b000
			}
		}

		if ch == 'c' {
			if matchRpc == 0b011 {
				matchRpc = 0b111
				continue
			} else {
				matchRpc = 0b000
			}
		}

	}

	return
}
