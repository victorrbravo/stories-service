package auth

const SuccessVerifyEmail = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Use2 app - Verify Email</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-KK94CHFLLe+nY2dmCWGMq91rCGa5gtU4mk92HdvYe+M/SXH301p5ILy+dN9+nJOZ" crossorigin="anonymous">
    <script src="https://unpkg.com/htmx.org@1.9.2" integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h" crossorigin="anonymous"></script>

</head>
<body class="container">

    <div class="row mt-4 g-4">
        <div class="col-8">
            <h1 class="mb-4">Activated successfully!</h1>
            <ul class="list-group fs-5 me-5" id="film-list">
                    <li class="list-group-item bg-primary text-white">{{.Name}} activated</li>
            </ul>
        </div>
    </div>
</body>
</html>
`

const AlreadyVerifiedEmail = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Use2 app - Verify Email</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-KK94CHFLLe+nY2dmCWGMq91rCGa5gtU4mk92HdvYe+M/SXH301p5ILy+dN9+nJOZ" crossorigin="anonymous">
    <script src="https://unpkg.com/htmx.org@1.9.2" integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h" crossorigin="anonymous"></script>

</head>
<body class="container">

    <div class="row mt-4 g-4">
        <div class="col-8">
            <h1 class="mb-4">User information</h1>
            <ul class="list-group fs-5 me-5" id="film-list">
                    <li class="list-group-item bg-primary text-white">{{.Name}} has already activated</li>
            </ul>
        </div>
     </div>
</body>
</html>
`

const CodewasResendTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Use2 app - Resend code</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-KK94CHFLLe+nY2dmCWGMq91rCGa5gtU4mk92HdvYe+M/SXH301p5ILy+dN9+nJOZ" crossorigin="anonymous">
    <script src="https://unpkg.com/htmx.org@1.9.2" integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h" crossorigin="anonymous"></script>

</head>
<body class="container">

    <div class="row mt-4 g-4">
        <div class="col-8">
            <h1 class="mb-4">Use2 app - Resend code</h1>
            <ul class="list-group fs-5 me-5" id="film-list">
                    <li class="list-group-item bg-primary text-white">The code for {{.Name}}  re-sent recently</li>
            </ul>
        </div>
     </div>
</body>
</html>
`

const CannotLoadChangePass = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Use2 app - Change Password</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-KK94CHFLLe+nY2dmCWGMq91rCGa5gtU4mk92HdvYe+M/SXH301p5ILy+dN9+nJOZ" crossorigin="anonymous">
    <script src="https://unpkg.com/htmx.org@1.9.2" integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h" crossorigin="anonymous"></script>

</head>
<body class="container">

    <div class="row mt-4 g-4">
        <div class="col-8">
            <h1 class="mb-4">Use2 app - Change Password</h1>
            <ul class="list-group fs-5 me-5" id="film-list">
                    <li class="list-group-item bg-primary text-white">Can't load change pass for {{.Name}}</li>
            </ul>
        </div>
     </div>
</body>
</html>
`

const ExpiredTokenChangePass = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Use2 app - Change Password</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-KK94CHFLLe+nY2dmCWGMq91rCGa5gtU4mk92HdvYe+M/SXH301p5ILy+dN9+nJOZ" crossorigin="anonymous">
    <script src="https://unpkg.com/htmx.org@1.9.2" integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h" crossorigin="anonymous"></script>

</head>
<body class="container">

    <div class="row mt-4 g-4">
        <div class="col-8">
            <h1 class="mb-4">Use2 app - Oh! your token has expired</h1>
            <ul class="list-group fs-5 me-5" id="film-list">
                    <li class="list-group-item bg-primary text-white">Can't load change pass, the token has expired for {{.Name}}</li>
            </ul>
        </div>
     </div>
</body>
</html>
`

const ExpiredTokenVerifyEmail = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Use2 - Verify Email</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-KK94CHFLLe+nY2dmCWGMq91rCGa5gtU4mk92HdvYe+M/SXH301p5ILy+dN9+nJOZ" crossorigin="anonymous">
    <script src="https://unpkg.com/htmx.org@1.9.2" integrity="sha384-L6OqL9pRWyyFU3+/bjdSri+iIphTN/bvYyM37tICVyOJkWZLpP2vGn6VUEXgzg6h" crossorigin="anonymous"></script>

</head>
<body class="container">

    <div class="row mt-4 g-4">
        <div class="col-8">
            <h1 class="mb-4">Use2 app - Oh! your token has expired</h1>
            <ul class="list-group fs-5 me-5" id="film-list">
                    <li class="list-group-item bg-primary text-white">Token has Expired for {{.Name}} </li>
            </ul>
	    </div>
			<div class="row mt-4 g-4">
              <li>Would you like us to resend the verification code to your email?  click on this <a href="{{.LinkResend}}">link</a></li>
	</div>
    </div>
	
</body>
</html>

`
