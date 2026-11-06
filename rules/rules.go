//go:build ruleguard

package gorules

import "github.com/quasilyte/go-ruleguard/dsl"

// noElse reports if/else statements so control flow can use guard clauses.
func noElse(m dsl.Matcher) {
	m.Match(`if $_ { $*_ } else { $*_ }`).
		Report("avoid else: use a guard clause and return early")
}

// noCallInLoop reports common database and network calls inside loops.
func noCallInLoop(m dsl.Matcher) {
	m.Match(
		`for $_ { $*_; $_.$m($*_); $*_ }`,
		`for $_ := range $_ { $*_; $_.$m($*_); $*_ }`,
		`for $_, $_ := range $_ { $*_; $_.$m($*_); $*_ }`,
		`for $*_; $*_; $*_ { $*_; $_.$m($*_); $*_ }`,
	).
		Where(m["m"].Text.Matches(`^(Query|QueryRow|QueryContext|QueryRowContext|Exec|ExecContext|Get|Post|Do|Set)$`)).
		Report("[error] DB/network call inside loop: batch it or move it out")
}

// preferErrName reports error variables whose names end in Err or err.
func preferErrName(m dsl.Matcher) {
	m.Match(
		`$name := $expr`,
		`$name, $_ := $expr`,
		`$_, $name := $expr`,
		`$name = $expr`,
		`$name, $_ = $expr`,
		`$_, $name = $expr`,
	).
		Where(m["name"].Text.Matches(`(?i)^.+err$`) && m["name"].Type.AssignableTo(`error`)).
		Report("[error] prefer err instead of a name ending in err")
}

func reviewFunction(m dsl.Matcher) {
	m.Match(
		`func $name($*_) $*_ { $*_ }`,
		`func ($_ $_) $name($*_) $*_ { $*_ }`,
	).
		Report("[review] Review this newly added function:\n" +
		    "1. Naming: Prefer a verb-first name when the function performs an action.\n" +
        	"2. Reusability: Consider whether the function can be reused and should be extracted into util package",
		)
}

func reviewSwitchStatement(m dsl.Matcher) {
	m.Match(
		`switch $x { case $first: $*_; case $second: $*_ }`,
		`switch $x { case $value: $*_; default: $*_ }`,
		`switch { case $first: $*_; case $second: $*_ }`,
		`switch { case $value: $*_; default: $*_ }`,
		`switch $x.(type) { case $first: $*_; case $second: $*_ }`,
		`switch $x.(type) { case $value: $*_; default: $*_ }`,
		`switch $name := $x.(type) { case $first: $*_; case $second: $*_ }`,
		`switch $name := $x.(type) { case $value: $*_; default: $*_ }`,
	).
		Report("[review] Review this two-scenario switch:\n" +
        "1. Prefer switch only when there are more than two scenarios.\n" +
        "2. Ensure each case condition, bodies is expressed through descriptive helper functions.\n",)
}

func boolReturns(m dsl.Matcher) {
	m.Match(
		`func $name($*_) ($*_) { $*_ }`,
		`func $name($*_) bool { $*_ }`,
		`func ($_ $_) $name($*_) ($*_) { $*_ }`,
		`func ($_ $_) $name($*_) bool { $*_ }`,
	).
		Where(
			!m["name"].Text.Matches(`^(Is|Has|Can|Should|Contains|Exists)`) &&
				m["name"].Type.Is(`func($*_) ($*_, bool, $*_)`),
		).
		Report("[review] This bool return may represent success/failure. Consider whether returning an error would be more expressive.")
}
