package auth

import "fmt"

func CreateEmailForUserAuth(username string, code string) string {
	filledTemplate := fmt.Sprintf(`
<html>
<head>
<title>GoAlign Authorization Code</title>  
</head>

<body>
<p>Hello %s,</p>

<p>Your GoAlign authorization code is: <b>%s</b></p> 

<p>Please enter this code in the GoAlign app to authenticate your account.</p>

<p>This code is valid for the next 10 minutes.</p>

<p>Thank you for using GoAlign!</p>

<p>The GoAlign Team</p>
</body>
</html>
  `, username, code)

	return filledTemplate
}

func CreateNonHtmlStringForUserAuth(username string, code string) string {
	filledTemplate := fmt.Sprintf(`Hello %s, Your authorization code is %s. Please enter this code in the GoAlign app to authenticate your account. This code is valid for the next 10 minutes. Thank you for using GoAlign! The GoAlign Team`, username, code)
	return filledTemplate
}
