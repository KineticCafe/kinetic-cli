package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
)

const (
	errorRed  = lipgloss.Color("#cc4444")
	purple    = lipgloss.Color("#7700dd")
	gray      = lipgloss.Color("#999999")
	lightGray = lipgloss.Color("#dddddd")
)

type resultStyles struct {
	header, rowOdd, rowEven, borderStyle lipgloss.Style
	borderType                           lipgloss.Border
	width                                int
}

type errorStyles struct {
	message lipgloss.Style
}

type styles struct {
	result resultStyles
	error  errorStyles
}

type styleOption func(*styles)

func newReportStyles(c *Config, options ...styleOption) styles {
	r := lipgloss.NewRenderer(c.stdout)
	cellStyle := r.NewStyle().Padding(0, 1)

	styles := styles{
		result: resultStyles{
			header:      r.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center),
			rowOdd:      cellStyle.Foreground(gray),
			rowEven:     cellStyle.Foreground(lightGray),
			borderStyle: r.NewStyle().Foreground(purple),
			borderType:  lipgloss.RoundedBorder(),
		},
		error: errorStyles{
			message: r.NewStyle().Foreground(errorRed).Bold(true),
		},
	}

	for _, option := range options {
		if option != nil {
			option(&styles)
		}
	}

	return styles
}

func withWidth(width int) styleOption {
	return func(styles *styles) {
		styles.result.width = width
	}
}

func reportResults(headers []string, results [][]string, styles styles) {
	if len(results) < 1 {
		return
	}

	tbl := table.New().
		Border(styles.result.borderType).
		BorderStyle(styles.result.borderStyle).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {

			case row == 0:
				return styles.result.header
			case row%2 == 0:
				return styles.result.rowEven
			default:
				return styles.result.rowOdd
			}
		}).
		Headers(headers...).
		Rows(results...)

	if styles.result.width > 0 {
		tbl = tbl.Width(styles.result.width)
	}

	fmt.Println(tbl)
}

func reportErrors(errors [][]string, styles styles) {
	if len(errors) < 1 {
		return
	}

	fmt.Print(styles.error.message.Render("Errors:") + "\n\n")

	errs := list.New()

	for _, item := range errors {
		key := item[0]
		msg := list.New(item[1]).
			ItemStyle(styles.error.message).
			Enumerator(func(_ list.Items, _ int) string { return "" })

		errs.Items(key, msg)
	}

	fmt.Println(errs)
}
