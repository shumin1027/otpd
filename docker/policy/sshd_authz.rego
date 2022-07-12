# This package path should be passed with the authz_endpoint flag
# in the sshd PAM configuration file.
package sshd.authz

import input.display_responses
import input.pull_responses
import input.sysinfo

# By default, users are not authorized.
default allow = false

allow {
	# OTP Validate
	url := sprintf("http://otpd:18181/validate?name=%s&passcode=%s", [sysinfo.pam_username, display_responses.passcode])
	response := http.send({"method": "get", "url": url})
	response.body.inventory == true
}

errors["You cannot pass!"] {
	not allow
}
