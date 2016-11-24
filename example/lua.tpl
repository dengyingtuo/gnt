local t = {
    {{- $N := len .Rows | dec}}
    {{- range $i, $row := .Rows}}
    [{{.EscapePk}}] = { 
        {{- $M := len .Data | dec}}
        {{- range $j, $cell := .Data }} 
	    {{if .Conv -}}
                {{- $keys := .ExpandKeys}}
                {{- $vals := .ExpandValues}}
	        {{- $maxRowIdx := len $vals | dec}}
                {{- .Name}} = {
		    {{- range $irow, $row := $vals}} 
	            {{- $maxColIdx := len $row| dec}}
		    {{- if $keys}}
		    { {{- range $iv, $val := $row}}{{- index $keys $iv}}={{$val}} {{- if lt $iv $maxColIdx}}, {{end}} {{- end }}}
		    {{- else}}
	            	{{- $maxColIdx := len $row| dec}}
		    { {{- range $iv, $val := $row}}{{- $val}} {{- if lt $iv $maxColIdx}},{{end}} {{- end }}}
		    {{- end}}
                    {{- if lt $irow $maxRowIdx}},{{end -}}
                    {{- end}}
        }	
	    {{- else}}
                {{- .Name}}={{.EscapeValue -}} 
            {{end}}
            {{- if lt $j $M}},{{end -}}
        {{- end}}
    }
    {{- if lt $i $N}},{{end}}
    {{- end}}
}

require'metadata'.new((...), t)
