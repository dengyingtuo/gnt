local t = {
    {{- $N := len .Rows | dec}}
    {{- range $i, $row := .Rows}}
    [{{.EscapePk}}] = { 
        {{- $M := len .Data | dec}}
        {{- range $j, $cell := .Data }} 
	    {{if .Split1Vals -}}
	        {{- $maxColIdx := len .Split1Vals | dec}}
		{{- if ge $maxColIdx 0}}
                {{- .Name}} = { {{- range $iv, $val := .Split1Vals}}{{- $val}} {{- if lt $iv $maxColIdx}}, {{end}} {{- end}}}
                {{- if lt $j $M}},{{end -}}
                {{- end}}
	    {{- else if .Split2Vals}}
	        {{- $maxRowIdx := len .Split2Vals | dec}}
		{{- $keys := .SplitKeys}}
                {{- if ge $maxRowIdx 0}}
                {{- .Name}} = {
		    {{- range $irow, $row := .Split2Vals}} 
	            {{- $maxColIdx := len $row| dec}}
		    {{- if $keys}}
		    { {{- range $iv, $val := $row}}{{- index $keys $iv}}={{$val}} {{- if lt $iv $maxColIdx}}, {{end}} {{- end }}}
		    {{- else}}
	            {{- $maxColIdx := len $row| dec}}
		    { {{- range $iv, $val := $row}}{{- $val}} {{- if lt $iv $maxColIdx}}, {{end}} {{- end }}}
		    {{- end}}
                    {{- if lt $irow $maxRowIdx}},{{end -}}
                    {{- end}}
		}
                {{- end}}
                {{- if lt $j $M}},{{end -}}
		{{- else}}
		{{- $val := .EscapeValue}}
                {{- if $val}}
                {{- .Name}}={{$val -}} 
                {{- if lt $j $M}},{{end -}}
                {{- end}}
            {{- end}}
        {{- end}}
    }
    {{- if lt $i $N}},{{- end}}
    {{- end}}
}

require'metadata'.new((...), t)
