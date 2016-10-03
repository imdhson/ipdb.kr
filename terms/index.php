<!DOCTYPE HTML>
<!--
	Stellar by HTML5 UP
	html5up.net | @ajlkn
	Free for personal and commercial use under the CCA 3.0 license (html5up.net/license)
-->
<?php
	$ip=$_SERVER['REMOTE_ADDR'];
 ?>
<html>
	<head>
		<?php
			echo "<title>$ip</title>";
		 ?>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
		<!--[if lte IE 8]><script src="assets/js/ie/html5shiv.js"></script><![endif]-->
		<link rel="stylesheet" href="assets/css/main.css" />
		<!--[if lte IE 9]><link rel="stylesheet" href="assets/css/ie9.css" /><![endif]-->
		<!--[if lte IE 8]><link rel="stylesheet" href="assets/css/ie8.css" /><![endif]-->
	</head>
	<body>

		<!-- Wrapper -->
			<div id="wrapper">

				<!-- Header -->
					<header id="header" class="alt">
						<?php
							echo "<h1>$ip</h1>";
						 ?>
						<h2>Your IP Address.</h2>
						<p>by <a href="http://ipdb.kr">ipdb.kr</a></p>
					</header>

				<!-- Nav -->
					<nav id="nav">
						<ul>
							<li><a href="#terms" class="active">Terms</a></li>
						</ul>
					</nav>

				<!-- Main -->
					<div id="main">

						<!-- Introduction -->
							<section id="terms" class="main">
								<div class="spotlight">
									<div class="content">
										<header class="major">
											<h2>IP Dashboard Terms of use.</h2>
										</header>
										<p>Accepting here, Your IP is record to "IP Dashboard". And anyone can see. Ofcourse, you can erase it.</p>
										<p>여기에 동의하면, "IP Dashboard"에 당신의 IP가 기록되고 모두가 볼 수 있습니다. 물론 지울 수 있습니다.</p>
										<ul class="actions">
											<li><a href="sendaccept" class="button">Accept</a></li>
											<p>If you can't proceed: <a href="http://ipdb.kr/accept">Click here</a>.</p>
											<p>진행할 수 없다면: <a href="http://ipdb.kr/accept">여기를 누르세요</a>.</p>
											<li><a href="sendreject" class="button">Reject</a></li>
										</ul>
									</div>
								</div>
							</section>
				<!-- Footer -->
					<footer id="footer">
						<p class="copyright">&copy; <a href="http://didlast.com/">dhson</a>. Design: <a href="https://html5up.net">HTML5 UP</a>. <a href="http://ipdb.kr/terms">Terms</a>. <a href="#">Go to Top</a>.</p>
					</footer>
			</div>

		<!-- Scripts -->
			<script src="assets/js/jquery.min.js"></script>
			<script src="assets/js/jquery.scrollex.min.js"></script>
			<script src="assets/js/jquery.scrolly.min.js"></script>
			<script src="assets/js/skel.min.js"></script>
			<script src="assets/js/util.js"></script>
			<!--[if lte IE 8]><script src="assets/js/ie/respond.min.js"></script><![endif]-->
			<script src="assets/js/main.js"></script>
	</body>
</html>
