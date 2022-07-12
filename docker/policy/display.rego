# The package path. This should be passed with the display_endpoint flag
# in the PAM configuration file.
package display

display_spec = [
	{
		"message": "Welcome to the OPA-PAM demonstration.",
		"style": "info",
	},
	{
		"message": "Please enter your OTP passcode: ",
		"style": "prompt_echo_on",
		"key": "passcode",
	}
]
