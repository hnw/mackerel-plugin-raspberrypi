# mackerel-plugin-raspberrypi

Raspberry PiのSoC温度・動作クロック・動作電圧・スロットル状態を取得するMackerelプラグイン

## 概要

Raspberry Pi の `vcgencmd` コマンドを利用して、SoC温度・動作クロック・動作電圧・スロットル状態を取得するMackerelプラグインです。

取得できる情報は以下の通りです。

- SoC温度 [℃]
- 動作クロック [MHz]
  - ARMコア
  - VC4スカラーコア
  - その他
- 動作電圧 [V]
  - VC4コア電圧
  - SDRAMコア電圧
  - その他
- スロットル状態
  - 低電圧を検出
  - （CPU高温のため[^1]）ARMの動作クロックを制限している
  - （低電圧のため[^1]）クロックダウンしている
  - CPU温度のソフトリミット到達
- スロットル状態の履歴
  - 低電圧を検出
  - （CPU高温のため[^1]）ARMの動作クロックを制限した
  - （低電圧のため[^1]）クロックダウンした
  - CPU温度のソフトリミット到達

[^1]: 参考: [Capping vs throttling \- what's the difference?](https://forums.raspberrypi.com/viewtopic.php?t=276404)

## 利用例

![temperature](https://raw.githubusercontent.com/hnw/mackerel-plugin-raspberrypi/images/temperature.png)

![clock](https://raw.githubusercontent.com/hnw/mackerel-plugin-raspberrypi/images/clock.png)

## 対応アーキテクチャ

GitHub Releasesで以下のアーキテクチャのバイナリを提供しています。

- arm(ARMv6)
- arm64

ビルドすれば他のアーキテクチャ用のバイナリも作れますが、無意味だと思います。

## インストール

`mkr plugin install`に対応しているので`mkr`を使うのがお勧めです。

```
$ sudo mkr plugin install --upgrade hnw/mackerel-plugin-raspberrypi
```

## Mackerelに登録する

`mackerel-agent.conf` を編集して `mackerel-agent` を再起動してください。しばらく待っているとカスタムメトリックが増えているはずです。

```
[plugin.metrics.raspberrypi]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-raspberrypi"
```
