/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package login

var (
	HTMLAuthSuccessful = `<!DOCTYPE html>
<html>
	<head>
		<title>Chainguard Authentication</title>
		<link rel="icon" href="https://console.enforce.dev/favicon.ico"/>
		<style>
			:root { color-scheme: light dark; font-family: "Helvetica Neue", "Arial", sans-serif; height: 100%; }
			body { display: flex; justify-content: center; height: 100%; margin: 0 10%; }
			.container { display: flex; justify-content: center; flex-direction: column; gap: 20px; }
			.title { font-size: 3em; }
			.content { font-size: 2em; }
			.links { display: flex; justify-content: center; font-size: 1.1em; margin-top: 10px; gap: 40px; user-select: none; }
			.link { color: #fff; text-decoration: none; padding: 16px; border: 1px solid #444CE4; border-radius: 8px; flex: 1 1 0; text-align: center; }
			.link:hover { opacity: 0.8; }
			.header { position: absolute; top: 30px; left: 20px; font-size: 25px; display: flex; height: 40px; line-height: 40px; gap: 5px; text-decoration: none; color: inherit; }
			.header>img { width: 40px; }
			.logo { color: #6363ec; font-weight: bold; }
			.chainctl { color: #6363ec; font-family: "Futura", "Consolas", sans-serif }
			.login { background: #fff; color: #111; }
			.offline .if-online, .if-offline { display: none; }
			.offline .if-offline { display: inline; }
			.login { background: #444CE4; color: #fff; }
			.login>.logo { color: #fff; }
			@media (prefers-color-scheme: light) {
				.logo { color: #4445e7 }
				.link { color: #111; }
				.login { color: #fff; }
				.link:hover { opacity: 0.7; }
				.login:hover { opacity: 0.9; }
			}
		</style>
	</head>
	<body>
		<script>
			if (!(window && window.navigator && window.navigator.onLine)) {
				document.body.className = "offline";
			}
		</script>
		<div class="header if-offline logo">Chainguard</div>
		<a class="header if-online" href="https://chainguard.dev">
			<img src="https://console.enforce.dev/logo512.png" onerror="this.remove()" />
			<span class="logo">Chainguard</span>
		</a>
		<div class="container">
			<div class="title">
				<span class="logo">Chainguard</span>
				<span> authentication successful!</span>
			</div>
			<div class="content">
				<span>You may now close </span>
				<span class="if-online">this page or visit:</span>
				<span class="if-offline">this page.</span>
			</div>
			<div class="links if-online">
				<a href="https://console.enforce.dev" class="link login"><span class="logo">Chainguard Console</span></a>
				<a href="https://edu.chainguard.dev" class="link"><span class="logo">Chainguard Academy</span></a>
				<a href="https://www.chainguard.dev/unchained" class="link"><span class="logo">Chainguard Blog</span></a>
			</div>
		</div>
	</body>
</html>
`
)
