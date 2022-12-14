import 'dart:convert';
import 'dart:math';

import 'package:client/manage/event.dart';
import 'package:client/model/event.dart';
import 'package:client/tool/date.dart';
import 'package:crypto/crypto.dart' as crypto;
import 'package:flutter/material.dart';
import 'package:package_info_plus/package_info_plus.dart';

/*
* 工具方法
* @author wuxubaiyang
* @Time 2022/9/8 15:09
*/
class Tool {
  // 生成id
  static String genID({int? seed}) {
    var time = DateTime.now().millisecondsSinceEpoch;
    return md5("${time}_${Random(seed ?? time).nextDouble()}");
  }

  // 生成时间戳签名
  static String genDateSign() =>
      DateTool.formatDate(DatePattern.dateSign, DateTime.now());

  // 计算md5
  static String md5(String value) =>
      crypto.md5.convert(utf8.encode(value)).toString();

  // 获取屏幕宽度
  static double getScreenWith(BuildContext context) =>
      MediaQuery.of(context).size.width;

  // 获取屏幕高度
  static double getScreenHeight(BuildContext context) =>
      MediaQuery.of(context).size.height;

  // 获取状态栏高度
  static double getStatusBarHeight(BuildContext context) =>
      MediaQuery.of(context).padding.top;

  // 获取应用主题色
  static Color getAccentColor(BuildContext context) =>
      Theme.of(context).colorScheme.primary;

  // 获取应用名
  static Future<String> get appName async {
    var packageInfo = await PackageInfo.fromPlatform();
    return packageInfo.appName;
  }

  // 获取应用包名
  static Future<String> get packageName async {
    var packageInfo = await PackageInfo.fromPlatform();
    return packageInfo.packageName;
  }

  // 获取版本号
  static Future<String> get buildNumber async {
    var packageInfo = await PackageInfo.fromPlatform();
    return packageInfo.buildNumber;
  }

  // 获取版本名
  static Future<String> get version async {
    var packageInfo = await PackageInfo.fromPlatform();
    return packageInfo.version;
  }
}

// 更新全局样式
void updateGlobalTheme(ThemeData themeData) =>
    eventManage.send(ThemeEvent(themeData: themeData));

// 将map转为query的url
String toQueryUrl(String url, Map<String, dynamic> params) {
  if (params.isNotEmpty) url += "?";
  for (var entry in params.entries) {
    url += "&${entry.key}=${entry.value}";
  }
  return url;
}

// map深层路径索引
V? findInMap<V>(Map map, String path) {
  if (path.isEmpty) return map as V;
  var paths = path.split(".");
  dynamic temp = map;
  for (var it in paths) {
    temp = temp[it];
    if (null == temp) return null;
  }
  return temp as V;
}
