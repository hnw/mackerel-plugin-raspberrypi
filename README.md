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
  - （高温のため）ARMの動作クロックに制限がかかっている
  - （低電圧のため）クロックダウンしている
  - CPU温度のソフトリミット到達

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

```
[plugin.metrics.raspberrypi]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-raspberrypi"
```
