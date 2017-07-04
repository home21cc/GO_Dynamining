{{define "base"}}
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

<!-- Begin Page content -->

<!-- Fixed navbar -->
<nav class="navbar navbar-default navbar-fixed-top">
    <div class="container">
        <div class="navbar-header">
            <img src="static/help.jpg" width="48" height="48"/>
        </div>
        <div id="navbar"  class="collapse navbar-collapse">
            <ul class="nav nav-pills">
                <li class="active"><a href="/Information">Information</a></li>
                <li><a href="/Rawmaterial">Raw Material</a></li>
                <li><a href="/Word">Word</a></li>
                <li><a href="/Sentence">Sentence</a></li>
                <li><a href="/Application">Applocation</a></li>
                <li><a href="/About">About</a></li>
                <li><a href="/Logout">Logout</a></li>
            </ul>
        </div><!--/.nav-collapse -->
    </div>
</nav>

{{template "body" .}}

<footer class="footer">
    <div class="container">
        <p class="text-muted">Create by dm,Park. </p>
    </div>
</footer>
</body>
</html>
{{end}}

