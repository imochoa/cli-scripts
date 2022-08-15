package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func wkhtmltopdfVersion() string {

	out, err := exec.Command("wkhtmltopdf", "--version").Output()
	stdout := string(out)
	if err != nil || len(stdout) <= 0 {
		panic("wkhtmltopdf not found! Install with 'sudo apt install wkhtmltopdf'")
	}
	version := strings.Join(strings.Fields(stdout)[1:], " ")
	return version

}

func pdfinfoVersion() string {

	out, err := exec.Command("pdfinfo", "-v").CombinedOutput()
	stdout := string(out)
	if err != nil || len(stdout) <= 0 {
		panic("pdfinfo not found! Install with 'sudo apt install pdfinfo'")
	}
	version := strings.Join(strings.Fields(stdout)[2:3], " ")

	return version
}

func getPdfPageCount(path string) int {

	out, err := exec.Command("pdfinfo", path).CombinedOutput()
	stdout := string(out)
	if err != nil || len(stdout) <= 0 {
		panic("pdfinfo failed!")
	}
	fields := strings.Fields(stdout)
	pageCount, _ := strconv.Atoi(fields[31])
	return pageCount

}

func HTML2LongPDF(inHTML, outPDF string, H, W int, units string) {

	//    # wkhtmltopdf -T 0 -B 0 --page-width "${W}mm" --page-height "${H}mm" "$1" "${AUX_PDF}"
	//
	//    # -B, --margin-bottom <unitreal>      Set the page bottom margin
	//    # -L, --margin-left <unitreal>        Set the page left margin (default 10mm)
	//    # -R, --margin-right <unitreal>       Set the page right margin (default 10mm)
	//    # -T, --margin-top <unitreal>         Set the page top margin
	//
	//    # -d, --dpi <dpi>                     Change the dpi explicitly (this has no
	//    #                                     effect on X11 based systems) (default 96)
	//
	//    # --background                        Do print background (default)

	// Initial estimate
	// tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("%d.pdf", rand.Float64()))
	// TODO get a good random path!
	tmpPath := filepath.Join(os.TempDir(), "temppath.pdf")

	// How many pages of (W,H) are required?
	_, _ = exec.Command("wkhtmltopdf",
		"-T", "0", "-B", "0",
		"--page-width", fmt.Sprintf("%s%s", strconv.Itoa(W), units),
		"--page-height", fmt.Sprintf("%s%s", strconv.Itoa(H), units),
		inHTML, tmpPath,
	).Output()

	pageCount := getPdfPageCount(tmpPath)

	// Make long PDF
	_, _ = exec.Command("wkhtmltopdf",
		"-T", "0", "-B", "0",
		"--page-width", fmt.Sprintf("%s%s", strconv.Itoa(W), units),
		"--page-height", fmt.Sprintf("%s%s", strconv.Itoa(H*pageCount), units),
		inHTML, outPDF,
	).Output()

}

func main() {

	inHTML := flag.String("input", "", "Input HTML")
	outPDF := flag.String("output", "", "Output PDF")
	flag.Parse()

	if *inHTML == "" || *outPDF == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Inputs
	*inHTML, _ = filepath.Abs(*inHTML)
	*outPDF, _ = filepath.Abs(*outPDF)

	// Check dependencies
	log.Printf("Found [wkhtmltopdf] (Version: %s)", wkhtmltopdfVersion())
	log.Printf("Found [pdfinfo] (Version: %s)", pdfinfoVersion())

	// A4  210 x 297 mm
	var A4H, A4W int = 297, 210
	A4U := "mm"

	HTML2LongPDF(*inHTML, *outPDF, A4H, A4W, A4U)

}
