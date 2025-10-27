# Wails App Launcher

This application, built with Wails, provides a simple desktop interface for launching various scripts from a remote server.

## Features

-   **Remote Script Execution:** Fetches a list of available scripts from a remote server and executes them on the user's machine.
-   **Multi-language Support:** Supports running scripts written in Python, Ruby, C, and shell commands.
-   **Dynamic Configuration:** Allows users to easily change the server address from which to fetch scripts.
-   **Cross-Platform:** Works on Windows, macOS, and Linux.

## How It Works

The application starts by fetching a `manifest.json` file from a configured server address. This JSON file contains a list of scripts that can be executed. Each entry in the manifest includes:

-   `name`: The display name of the script.
-   `description`: A brief description of what the script does.
-   `filename`: The name of the script file to download or the shell command to execute.
-   `language`: The language of the script (e.g., `python`, `ruby`, `c`, `shell`).

When a user chooses to run a script:

1.  The application downloads the script file from the server (unless it's a direct shell command).
2.  For C code, the application first compiles it into an executable using `gcc`.
3.  A new terminal window is opened, and the script or compiled binary is executed.

## Getting Started

### Prerequisites

-   Go (for the backend)
-   Node.js/npm (for the frontend)
-   Wails CLI: Follow the installation instructions at https://wails.io/docs/gettingstarted/installation.
-   For running C scripts, `gcc` must be installed and available in your system's PATH.

### Setting Up the Server

This application requires a simple HTTP server to provide the script files and the `manifest.json`. A simple way to do this is to use Python's built-in HTTP server.

1.  Navigate to the `server-files` directory in this project:
    ```sh
    cd server-files
    ```

2.  Start a local HTTP server. For Python 3, you can run:
    ```sh
    python -m http.server 8080
    ```
    This will serve the files in the `server-files` directory on `http://localhost:8080`.

### Running the Application

1.  In the project's root directory, run the application in development mode:
    ```sh
    wails dev
    ```

2.  The application window will open. The default server address is `http://localhost:8080/`. If you are using a different address for your server, you can change it in the application's UI.

3.  You should see a list of scripts from the `manifest.json` file. Click "Run" to execute a script.

## Building

To build a redistributable, production-ready package, use the following command in the project's root directory:

```sh
wails build
```

This will generate a native application for your platform in the `build/bin` directory.

---

# Wails App Launcher (日本語)

このアプリケーションはWailsで構築されており、リモートサーバーから様々なスクリプトを起動するためのシンプルなデスクトップインターフェースを提供します。

## 機能

-   **リモートスクリプト実行:** リモートサーバーから利用可能なスクリプトのリストを取得し、ユーザーのマシンで実行します。
-   **多言語サポート:** Python, Ruby, C, シェルコマンドで書かれたスクリプトの実行をサポートします。
-   **動的設定:** スクリプトを取得するサーバーアドレスをユーザーが簡単に変更できます。
-   **クロスプラットフォーム:** Windows, macOS, Linuxで動作します。

## 仕組み

アプリケーションは、設定されたサーバーアドレスから`manifest.json`ファイルを取得することから始まります。このJSONファイルには、実行可能なスクリプトのリストが含まれています。マニフェストの各エントリには、次のものが含まれます。

-   `name`: スクリプトの表示名。
-   `description`: スクリプトの動作に関する簡単な説明。
-   `filename`: ダウンロードするスクリプトファイルの名前、または実行するシェルコマンド。
-   `language`: スクリプトの言語（例: `python`, `ruby`, `c`, `shell`）。

ユーザーがスクリプトの実行を選択すると：

1.  アプリケーションはサーバーからスクリプトファイルをダウンロードします（直接のシェルコマンドでない場合）。
2.  Cコードの場合、アプリケーションはまず`gcc`を使用して実行可能ファイルにコンパイルします。
3.  新しいターミナルウィンドウが開き、スクリプトまたはコンパイルされたバイナリが実行されます。

## はじめに

### 前提条件

-   Go (バックエンド用)
-   Node.js/npm (フロントエンド用)
-   Wails CLI: https://wails.io/docs/gettingstarted/installation のインストール手順に従ってください。
-   Cスクリプトを実行するには、`gcc`がインストールされ、システムのPATHで利用可能である必要があります。

### サーバーのセットアップ

このアプリケーションは、スクリプトファイルと`manifest.json`を提供するためのシンプルなHTTPサーバーを必要とします。これを行う簡単な方法は、Pythonの組み込みHTTPサーバーを使用することです。

1.  このプロジェクトの`server-files`ディレクトリに移動します:
    ```sh
    cd server-files
    ```

2.  ローカルHTTPサーバーを起動します。Python 3の場合は、次のように実行できます:
    ```sh
    python -m http.server 8080
    ```
    これにより、`server-files`ディレクトリ内のファイルが`http://localhost:8080`で提供されます。

### アプリケーションの実行

1.  プロジェクトのルートディレクトリで、開発モードでアプリケーションを実行します:
    ```sh
    wails dev
    ```

2.  アプリケーションウィンドウが開きます。デフォルトのサーバーアドレスは`http://localhost:8080/`です。サーバーに別のアドレスを使用している場合は、アプリケーションのUIで変更できます。

3.  `manifest.json`ファイルからのスクリプトのリストが表示されます。「実行」をクリックしてスクリプトを実行します。

## ビルド

再配布可能な、本番環境対応のパッケージをビルドするには、プロジェクトのルートディレクトリで次のコマンドを使用します:

```sh
wails build
```

これにより、`build/bin`ディレクトリにプラットフォーム用のネイティブアプリケーションが生成されます。
