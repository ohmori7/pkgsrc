package main

import "netbsd.org/pkglint/trace"

// SubstContext records the state of a block of variable assignments
// that make up a SUBST class (see `mk/subst.mk`).
type SubstContext struct {
	id            string
	stage         string
	message       string
	curr          *SubstContextStats
	inAllBranches SubstContextStats
	filterCmd     string
}

func NewSubstContext() *SubstContext {
	return &SubstContext{curr: &SubstContextStats{}}
}

type SubstContextStats struct {
	seenFiles     bool
	seenSed       bool
	seenVars      bool
	seenTransform bool
	prev          *SubstContextStats
}

func (st *SubstContextStats) Copy() *SubstContextStats {
	return &SubstContextStats{st.seenFiles, st.seenSed, st.seenVars, st.seenTransform, st}
}

func (st *SubstContextStats) And(other *SubstContextStats) {
	st.seenFiles = st.seenFiles && other.seenFiles
	st.seenSed = st.seenSed && other.seenSed
	st.seenVars = st.seenVars && other.seenVars
	st.seenTransform = st.seenTransform && other.seenTransform
}

func (st *SubstContextStats) Or(other SubstContextStats) {
	st.seenFiles = st.seenFiles || other.seenFiles
	st.seenSed = st.seenSed || other.seenSed
	st.seenVars = st.seenVars || other.seenVars
	st.seenTransform = st.seenTransform || other.seenTransform
}

func (ctx *SubstContext) Varassign(mkline MkLine) {
	if trace.Tracing {
		trace.Stepf("SubstContext.Varassign %#v %v#", ctx.curr, ctx.inAllBranches)
	}

	varname := mkline.Varname()
	varcanon := mkline.Varcanon()
	varparam := mkline.Varparam()
	op := mkline.Op()
	value := mkline.Value()
	if varcanon == "SUBST_CLASSES" || varcanon == "SUBST_CLASSES.*" {
		classes := fields(value)
		if len(classes) > 1 {
			mkline.Warnf("Please add only one class at a time to SUBST_CLASSES.")
		}
		if ctx.id != "" && ctx.id != classes[0] {
			complete := ctx.IsComplete()
			id := ctx.id
			ctx.Finish(mkline)
			if !complete {
				mkline.Warnf("Subst block %q should be finished before adding the next class to SUBST_CLASSES.", id)
			}
		}
		ctx.id = classes[0]
		return
	}

	switch varcanon {
	case "SUBST_STAGE.*":
	case "SUBST_MESSAGE.*":
	case "SUBST_FILES.*":
	case "SUBST_SED.*":
	case "SUBST_VARS.*":
	case "SUBST_FILTER_CMD.*":

	default:
		if ctx.id != "" {
			mkline.Warnf("Foreign variable %q in SUBST block.", varname)
		}
		return
	}

	if ctx.id == "" {
		mkline.Warnf("SUBST_CLASSES should come before the definition of %q.", varname)
		ctx.id = varparam
	}

	if varparam != ctx.id {
		if ctx.IsComplete() {
			// XXX: This code sometimes produces weird warnings. See
			// meta-pkgs/xorg/Makefile.common 1.41 for an example.
			ctx.Finish(mkline)

			// The following assignment prevents an additional warning,
			// but from a technically viewpoint, it is incorrect.
			ctx.id = varparam
		} else {
			mkline.Warnf("Variable %q does not match SUBST class %q.", varname, ctx.id)
			return
		}
	}

	switch varcanon {
	case "SUBST_STAGE.*":
		ctx.dupString(mkline, &ctx.stage, varname, value)
		if value == "pre-patch" || value == "post-patch" {
			fix := mkline.Autofix()
			fix.Warnf("Substitutions should not happen in the patch phase.")
			fix.Explain(
				"Performing substitutions during post-patch breaks tools such as",
				"mkpatches, making it very difficult to regenerate correct patches",
				"after making changes, and often leading to substituted string",
				"replacements being committed.",
				"",
				"Instead of pre-patch, use post-extract.",
				"Instead of post-patch, use pre-configure.")
			fix.Replace("pre-patch", "post-extract")
			fix.Replace("post-patch", "pre-configure")
			fix.Apply()
		}

		if G.Pkg != nil && (value == "pre-configure" || value == "post-configure") {
			if noConfigureLine := G.Pkg.vars.FirstDefinition("NO_CONFIGURE"); noConfigureLine != nil {
				mkline.Warnf("SUBST_STAGE %s has no effect when NO_CONFIGURE is set (in %s).",
					value, noConfigureLine.ReferenceFrom(mkline.Line))
				Explain(
					"To fix this properly, remove the definition of NO_CONFIGURE.")
			}
		}

	case "SUBST_MESSAGE.*":
		ctx.dupString(mkline, &ctx.message, varname, value)
	case "SUBST_FILES.*":
		ctx.dupBool(mkline, &ctx.curr.seenFiles, varname, op, value)
	case "SUBST_SED.*":
		ctx.dupBool(mkline, &ctx.curr.seenSed, varname, op, value)
		ctx.curr.seenTransform = true
	case "SUBST_VARS.*":
		ctx.dupBool(mkline, &ctx.curr.seenVars, varname, op, value)
		ctx.curr.seenTransform = true
	case "SUBST_FILTER_CMD.*":
		ctx.dupString(mkline, &ctx.filterCmd, varname, value)
		ctx.curr.seenTransform = true
	}
}

func (ctx *SubstContext) Directive(mkline MkLine) {
	if ctx.id == "" {
		return
	}

	if trace.Tracing {
		trace.Stepf("+ SubstContext.Directive %#v %v#", ctx.curr, ctx.inAllBranches)
	}
	dir := mkline.Directive()
	if dir == "if" {
		ctx.inAllBranches = SubstContextStats{true, true, true, true, nil}
	}
	if dir == "elif" || dir == "else" || dir == "endif" {
		if ctx.curr.prev != nil { // Don't crash on malformed input
			ctx.inAllBranches.And(ctx.curr)
			ctx.curr = ctx.curr.prev
		}
	}
	if dir == "if" || dir == "elif" || dir == "else" {
		ctx.curr = ctx.curr.Copy()
	}
	if dir == "endif" {
		ctx.curr.Or(ctx.inAllBranches)
	}
	if trace.Tracing {
		trace.Stepf("- SubstContext.Directive %#v %v#", ctx.curr, ctx.inAllBranches)
	}
}

func (ctx *SubstContext) IsComplete() bool {
	return ctx.id != "" &&
		ctx.stage != "" &&
		ctx.curr.seenFiles &&
		ctx.curr.seenTransform
}

func (ctx *SubstContext) Finish(mkline MkLine) {
	if ctx.id == "" {
		return
	}

	id := ctx.id
	if ctx.stage == "" {
		mkline.Warnf("Incomplete SUBST block: SUBST_STAGE.%s missing.", id)
	}
	if !ctx.curr.seenFiles {
		mkline.Warnf("Incomplete SUBST block: SUBST_FILES.%s missing.", id)
	}
	if !ctx.curr.seenTransform {
		mkline.Warnf("Incomplete SUBST block: SUBST_SED.%[1]s, SUBST_VARS.%[1]s or SUBST_FILTER_CMD.%[1]s missing.", id)
	}

	ctx.id = ""
	ctx.stage = ""
	ctx.message = ""
	ctx.curr = &SubstContextStats{}
	ctx.filterCmd = ""
}

func (ctx *SubstContext) dupString(mkline MkLine, pstr *string, varname, value string) {
	if *pstr != "" {
		mkline.Warnf("Duplicate definition of %q.", varname)
	}
	*pstr = value
}

func (ctx *SubstContext) dupBool(mkline MkLine, flag *bool, varname string, op MkOperator, value string) {
	if *flag && op != opAssignAppend {
		mkline.Warnf("All but the first %q lines should use the \"+=\" operator.", varname)
	}
	*flag = true
}
