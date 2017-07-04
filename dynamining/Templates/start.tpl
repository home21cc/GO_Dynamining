{{define "head"}}
    <title>Start</title>
{{end}}

{{define "body"}}
<div class="container">
	<div class="page-header">
		<h1>OGAM help system</h1>
	</div>
	<p class="lead">Do you need help from your computer ?</p>
</div>

<div class="container">

	<form action="/" class="form-signin" method="post">
		<h2 class="form-signin-heading">Please email</h2>
		<label for="inputEmail" class="sr-only">Email address</label>
		<input type="email" name="inputEmail" class="form-control" placeholder="Email address" required autofocus>
		<label for="inputPassword" class="sr-only">Password</label>
		<input type="password" name="inputPassword" class="form-control" placeholder="Password" >
        <p></p>
		<button class="btn btn-lg btn-primary btn-block" type="submit">email check in</button>
	</form>

</div> <!-- /container -->
{{end}}