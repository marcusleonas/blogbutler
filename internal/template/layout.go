package template

var LayoutTemplate = `{{ define "layout" }}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .PageTitle }} | {{ .SiteTitle }}</title>
</head>
<body>
    <main class="content">
        {{ template "content" . }}
    </main>
    <footer class="footer">
        <p>&copy; {{ .Copyright }}</p>
    </footer>
</body>
</html>
{{ end }}`
