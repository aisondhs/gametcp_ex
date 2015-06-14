<?php
error_reporting(0);
set_time_limit(0);

$host = "127.0.0.1";  
$port = 8989;
$socket = socket_create(AF_INET, SOCK_STREAM, SOL_TCP)or die("Could not create  socket\n");
$connection = socket_connect($socket, $host, $port) or die("Could not connet server\n");

$request = array("account"=>"ella","pwd"=>md5("123456"));
$reqContent = json_encode($request);
$msgId = 101;
$data = pack("n",strlen($reqContent)+4).pack("n", $msgId).$reqContent;
socket_write($socket, $data) or die("Write failed\n");

$rspdata = socket_read($socket, 1024);
    
$str = substr($rspdata,4)."\n";
$retData = json_decode($str,true);
print_r($retData);
echo "\n";

$request = array("verify"=>$retData["token"]);
$reqContent = json_encode($request);
$msgId = 0;
$data = pack("n",strlen($reqContent)+4).pack("n", $msgId).$reqContent;
socket_write($socket, $data) or die("Write failed\n");

$rspdata = socket_read($socket, 1024);
$str = substr($rspdata,4)."\n";
$retData = json_decode($str,true);
print_r($retData);
echo "\n";

socket_close($socket);
