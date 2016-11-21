{
    {{- $N := len .Rows | dec}}
    {{- range $i, $row := .Rows}}
    "{{.Pk}}" : { 
        {{- $M := len .Data | dec}}
        {{- range $j, $cell := .Data }} 
            "{{.Name}}": {{.EscapeValue}} 
            {{- if lt $j $M}},{{end}}
        {{- end}}
    }	
    {{- if lt $i $N}},{{end}}
    {{- end}}
}
