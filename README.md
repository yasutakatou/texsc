# texsc

自分でテキストコーパスをつくってみたい！けどネット記事やブログ、youtubeにネタが散在しててスクレイピングが難しい・・そんなあなたの**人力スクレイピング**を応援するツールです。

![demo](https://github.com/yasutakatou/texsc/blob/pic/demo.gif)

## 事前準備

おもちのWindowsにmecabとfasttextをインストールし、パスを通してください。<br>
[ここ](https://github.com/ikegami-yukino/mecab)やら[ここ](https://github.com/xiamx/fastText)からダウンロード出来るでしょう。


## 動かし方

以下のようにコンパイルするか、[リリースページからバイナリを拾って](https://github.com/yasutakatou/texsc/releases)ください。Go言語なのでバイナリ単体で動きます。

```
git clone https://github.com/yasutakatou/texsc
cd texsc
go build texsc.go
```

起動オプションは以下です。

|オプション名|デフォルト値|説明|
|:---|:---|:---|
|debug|false|デバッグモードにするかどうか。内部出力が色々見えます|
|label|__label__1|fasttextのラベル名。後から変えられます|
|model|model|モデルのファイル名。|
|file|test.txt|出力するコーパスのファイル名。後から変えられます|
|regexp|[0-9][0-9]:[0-9][0-9]|正規表現で不要な部分を消しながらコーパスを作れます。よって指定は正規表現です|
|predict|fasttext.exe predict-prob #MODEL#.bin #DATA#|予測確認用コマンド。このままOS上で実行されます|
|learn|fasttext.exe supervised -input #DATA# -output #MODEL# -epoch 1000|学習用コマンド。こちらもこのままOS上で実行されます|
|predictChar|88|予測モードへの切り替えキーコード。デフォルトは88でX|
|learnChar|90|学習コマンド実施のキーコード。デフォルトは90でZ|

キーコードの指定は[こちらのサイト](http://shanabrian.com/web/javascript/keycode.php)から確認してください。<br>

> regexpオプションのデフォルトパラメータはyoutubeの文字起こしに使えます。

![regexp](https://github.com/yasutakatou/texsc/blob/pic/regexp.png)


## 使い方

 - デモのようにctrl+cでクリップボードにテキストがコピーされるとテキストのコーパスが作成されます。

```
__label__1, なんで こんな クド 過ぎる タイトル な の か は 前回 を ご 参照 ください 。 記事 として は 業務 経験 に 漬け込ん だ 自動 化 の ツール 作っ た ので 見 て
```

 - 予測モードに切り替えてからクリップボードにテキストがコピーされると、選択したテキストで予測されます。

```
>>>  annotation: false
 -- -- -- -- predict! -- -- -- --
__label__1, 0.504672
```

 - 学習コマンド実施のキーコードではモデルの作成が行われます。

```
>>>  learning..
Done!
```

## その他

起動後、>>> のプロンプトからターミナルコマンドを実行できます

|コマンド名|説明|
|:---|:---|
|setLabel|ラベル名を変えます。即有効なのでラベル名を切り替えながらスクレイピングできます|
|setFile|コーパスのファイル名を変えます。即有効なのでファイル名を切り替えながらスクレイピングできます|
|showConfig|今のコンフィグを見ます|

> 設定は以下のようにコマンドに引数を与えて設定します。

```
>>> setLabel __label__2
```
