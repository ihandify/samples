<?php

require_once __DIR__ . '/../services.php';

header('Content-Type: application/json');

if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
    http_response_code(405);

    echo json_encode([
        'error' => 'Method Not Allowed'
    ]);

    exit;
}

$input = json_decode(file_get_contents('php://input'), true);

$engines = $input['engines'] ?? null;
$expiresInSeconds = $input['expiresInSeconds'] ?? null;

$result = generateScopedPublicKey(
    $engines,
    (int)$expiresInSeconds
);

if ($result === null) {
    http_response_code(502);

    echo json_encode([
        'error' => 'Upstream service error'
    ]);

    exit;
}

echo json_encode($result);
