<?php

if ($_SERVER['REQUEST_URI'] === '/') {
    readfile(__DIR__ . '/static/demo.html');
    exit;
}
