// spiritgen is the CLI front-end. It dispatches to subcommand handlers:
//
//	spiritgen tablet  --input <xlsx> --output <pdf>
//	spiritgen lantern --input <xlsx> --output <pdf> --title <행사명>
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	args := os.Args[2:]
	switch os.Args[1] {
	case "tablet":
		runTablet(args)
	case "lantern":
		runLantern(args)
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "❌ 알 수 없는 명령어: %s\n\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: spiritgen <subcommand> [options]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Subcommands:")
	fmt.Fprintln(os.Stderr, "  tablet   영가위패 PDF 생성")
	fmt.Fprintln(os.Stderr, "  lantern  인등/연등 PDF 생성")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "각 서브커맨드의 옵션은 'spiritgen <subcommand> --help' 로 확인하세요.")
}
