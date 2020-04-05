package main

func index(port string) string {
	return `<html>
		<head>
			<title>God</title>
			<style>
				body {
					margin: 0;
				}
		
				div {
					text-align: center;
					font-size: 20px;
					padding: 10px;
					font-family: sans-serif;
					background-color: #ddd;
				}
		
				a {
					color: black;
					text-decoration: none;
					margin: 0 10px;
				}
		
				a:hover {
					color: blue;
				}
		
				iframe {
					width: 100%;
					height: calc(100% - 40px);
					border: 0;
				}
			</style>
		</head>
		<body>
			<div>
				<a href="http://0.0.0.0:3790/ui/flamegraph" target="frame">CPU</a>
				<a href="http://0.0.0.0:3791/ui/flamegraph" target="frame">Mem</a>
				<a href="http://0.0.0.0:3792/ui/flamegraph" target="frame">Block</a>
				<a href="http://0.0.0.0:3793/ui/flamegraph" target="frame">Mutex</a>
				<a href="http://0.0.0.0:3794" target="frame">Trace</a>
				<a href="http://0.0.0.0:` + port + `/metrics" target="frame">Metrics</a>
			</div>
		
			<iframe name="frame" src="http://0.0.0.0:3790/ui/flamegraph">
			</iframe>
		</body>
	</html>`
}
