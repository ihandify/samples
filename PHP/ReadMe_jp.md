# PHP版 iHandify サンプルアプリケーション

このサンプルアプリケーションは、iHandify APIをアプリケーションに統合する方法を示したデモです。

## 前提条件

- Docker
- Docker Compose
- 有効な iHandify サブスクリプション
- 有効なシークレットAPIキー（Secret API Key）

## はじめに

### 1. リポジトリのクローン

```bash
git clone https://github.com/ihandify/samples.git
cd samples/PHP
```

### 2. シークレットAPIキーの取得

1. iHandify ダッシュボードにサインインします。
2. **API Key**（APIキー）に移動します。
3. **API Key**（APIキー）をコピーします。

### 3. 環境変数の設定

`.env` ファイルを開き、プレースホルダーの値を実際のシークレットAPIキーに置き換えます。

```env
API_KEY="Paste_your_Secret_API_Key_here"
```

### 4. アプリケーションの起動

Docker Compose を使用してアプリケーションを起動します。

```bash
docker compose up -d --build
```

コンテナが正常に実行されているか確認します。

```bash
docker compose ps
```

### 5. デモアプリケーションを開く

ブラウザを開き、次のURLにアクセスします。

```text
http://localhost:8020
```

これで、サンプルアプリケーションを使用した手書き文字認識を試すことができます。

---

## 認証フロー

このサンプルは、推奨されるセキュリティアーキテクチャに準拠しています。

```text
フロントエンド
    ↓
サンプルバックエンド
    ↓ （シークレットAPIキーを使用）
iHandify 認証API
    ↓
スコープ付き公開APIキー（Scoped Public API Key）
    ↓
フロントエンド
    ↓
iHandify 認識API
```

1. フロントエンドがサンプルバックエンドにリクエストを送信します。
```javascript
    const response = await fetch(KEY_API_URL, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            expiresInSeconds: SCOPED_KEY_EXPIRES_IN
        })
    })
```
2. バックエンドはシークレットAPIキーを使用して、iHandifyから「スコープ付き公開APIキー」をリクエストします。
```php
    $ch = curl_init();

    curl_setopt_array($ch, [
        CURLOPT_URL => rtrim($apiUrl, '/') . '/plan/auth/generate-scoped-public-key',
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_POST => true,
        CURLOPT_POSTFIELDS => $payload,
        CURLOPT_HTTPHEADER => [
            'Content-Type: application/json',
            'X-Api-Key: ' . $apiKey
        ],
        CURLOPT_TIMEOUT => 30
    ]);

    $response = curl_exec($ch);
```
3. スコープ付き公開APIキーがフロントエンドに返されます。
```javascript
    const result = await response.json();
    scopedPublicKey = result.data.scopedPublicKey;
```
4. フロントエンドは、そのスコープ付き公開APIキーを使用して iHandify 認識APIを呼び出します。
```javascript
    const fd = new FormData();
    fd.append("input", JSON.stringify(pattern));

    const response = await fetch(url, {
        method: 'POST',
        body: fd,
        signal: AbortSignal.timeout(REQUEST_TIMEOUT_MS),
        headers: {
            'x-scoped-key': scopedPublicKey || ''
        }
    });
```

> [!CAUTION]
**フロントエンドアプリケーション（クライアント側）にシークレットAPIキーを決して露出させないでください。**

---

## アプリケーションの停止

```bash
docker compose down
```

## トラブルシューティング

### コンテナの起動に失敗する場合

コンテナのログを確認してください。

```bash
docker compose logs
```

### 認証に失敗する場合

以下を確認してください。

- シークレットAPIキーが有効であること。
- `.env`ファイルにAPIキーが正しく設定されていること。
- サブスクリプションが有効（アクティブ）であること。


### ポート 8020 がすでに使用されている場合

`docker-compose.yml` 内のポートマッピングを変更するか、現在ポート 8020 を使用している他のアプリケーションを停止してください。

---

## 関連リソース

- ドキュメント: https://ihandify.com/docs
- その他のサンプルアプリケーション: https://github.com/ihandify

ご質問やサポートが必要な場合は、iHandify チームまでお問い合わせください。
