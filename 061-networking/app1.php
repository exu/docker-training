<?php

$response = file_get_contents("http://app2/");

?>
<h2>Got response from http:/app2/ API</h2>
<code style="color:silver; background: black; padding:1em">
<?php
print_r($response);
?>
</code>
