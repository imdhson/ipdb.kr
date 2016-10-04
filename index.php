<!DOCTYPE HTML>
<!--
	Stellar by HTML5 UP
	html5up.net | @ajlkn
	Free for personal and commercial use under the CCA 3.0 license (html5up.net/license)
-->
<?php
  if($_COOKIE[terms] == "accept") {
	$ip=$_SERVER['REMOTE_ADDR'];
	$date = date("Y-m-d H:i:s");
	$fp = fopen("iphistory.txt", "a");
	fwrite($fp, "<tr>"."<td>$date</td>"."<td>$ip</td>"."<tr>");
	fclose($fp);
	$fp1 = fopen("iphistory.txt","rb");
	$contents = fread($fp1, filesize("iphistory.txt"));
	fclose($fp1);
  if(filesize("iphistory.txt") > 1024) {
    $fp2 = fopen("iphistory.txt", "w");
    fwrite($fp2, "");
    fclose($fp2);
  }
}
else	{
	echo "<meta http-equiv='refresh' content='0;url=terms' />";
}
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
  <style>
  table {
      border-collapse: collapse;
      width: 100%;
  }

  td, th {
      border: 1px solid #dddddd;
      text-align: left;
      padding: 8px;
  }

  tr:nth-child(even) {
      background-color: #dddddd;
  }
  </style>
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
							<li><a href="#dash" class="active">Dashboard</a></li>
              <li><a href="terms" class="active">Terms</a></li>
						</ul>
					</nav>

				<!-- Main -->
					<div id="main">

						<!-- Introduction -->
							<section id="dash" class="main">
								<div class="spotlight">
									<div class="content">
										<header class="major">
											<h2>IP Dashboard</h2>
										</header>
                    <p><a href="#footer" class="active">Go to Bottom.</a></p>
                    <table>
                      <tr>
                        <th>Date</th>
                        <th>IP</th>
                      </tr>
										<?php
											echo "<p>$contents</p>";
										 ?>
                     </table>
                     <ul class="actions">
											<li><a href="sendremove" class="button">remove</a></li>
										</ul>
									</div>
								</div>
							</section>
				<!-- Footer -->
					<footer id="footer">
						<p class="copyright">&copy; <a href="http://didlast.com/">dhson</a>. Design: <a href="https://html5up.net">HTML5 UP</a>. <a href="https://github.com/hhammer1234/ipdb.kr">Github</a>. <a href="terms">Terms</a>. <a href="#">Go to Top</a>.</p>
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
