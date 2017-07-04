{{define "basic"}}
<html >
<head>{{template "head" .}}
    <meta charset="utf-8">
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css" type="text/css">
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap-theme.min.css" type="text/css">
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js"></script>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.2/jquery.min.js"></script>
    <link href="static/start.css" rel="stylesheet">
</head>
<body>
{{template "body" .}}
<p/>
<p/>
<p/>
<footer class="footer">
    <div class="container">
        <p class="text-muted">Create by dm, Park</p>
    </div>
</footer>
</body>
</html>
{{end}}
