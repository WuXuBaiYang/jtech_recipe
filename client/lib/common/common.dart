import 'dart:io';

import 'package:client/main.dart';
import 'package:client/page/home/home.dart';
import 'package:client/page/login.dart';
import 'package:client/tool/file.dart';
import 'package:flutter/widgets.dart';

/*
* 静态资源/通用静态变量
* @author wuxubaiyang
* @Time 2022/9/8 14:54
*/
class Common {
  // api基础地址
  static String baseUrl = debugMode ? _baseUrlDev : _baseUrlRelease;

  // api开发地址
  static const String _baseUrlDev = "http://192.168.16.50:9527/api";

  // api正式地址
  static const String _baseUrlRelease = "https://$remoteHost:9527/api";

  // 远程域名
  // static const String remoteHost = "api.jtech.live";
  static const String remoteHost = "35.73.34.38";

  // 加密盐
  static const String salt = "p5bDTyO6McoYeH";
}

/*
* 路由路径静态变量
* @author wuxubaiyang
* @Time 2022/9/8 14:55
*/
class RoutePath {
  // 创建路由表
  static Map<String, WidgetBuilder> get routes => {
        home: (c) => const HomePage(),
        login: (c) => const LoginPage(),
      };

  // 首页
  static const String home = "/home";

  // 登录页
  static const String login = "/login";
}

/*
* oss服务配置信息
* @author wuxubaiyang
* @Time 2022/9/8 15:39
*/
class OSSConfig {
  // 地址
  static String endPoint = debugMode ? _endPointDev : _endPointRelease;

  // 开发版地址
  static const String _endPointDev = Common.remoteHost;

  // 正式版地址
  static const String _endPointRelease = "";

  // 端口号
  static const int port = 9000;

  // 是否启用ssl
  static const useSSL = false;

  // 授权key
  static const accessKey = r"QFIknqJRoZOJV7E2";

  // 密钥
  static const secretKey = r"SV46qm1zsumsyz870q6nXytHTLOoyxZp";
}

/*
* IM配置信息
* @author wuxubaiyang
* @Time 2022/9/9 19:17
*/
class IMConfig {
  // 平台
  static int get platformCode {
    if (Platform.isAndroid) return 2;
    if (Platform.isIOS) return 1;
    return 0;
  }

  // api地址
  static const String apiAddr = "http://${Common.remoteHost}:10002";

  // ws地址
  static const String wsAddr = "ws://${Common.remoteHost}:10003";

  // 数据存储路径
  static Future<String?> get dataDir => FileTool.getDirPath(
        FileDirPath.imDocumentPath,
        root: FileDir.applicationDocuments,
      );

  // oss类型
  static const String objectStorage = "minio";

  // 日志等级
  static int get logLevel => debugMode ? 6 : 1;
}

/*
* 目录管理
* @author wuxubaiyang
* @Time 2022/9/9 17:57
*/
class FileDirPath {
  // 图片缓存路径
  static const String imageCachePath = "imageCache";

  // 视频缓存路径
  static const String videoCachePath = "videoCache";

  // 音频缓存路径
  static const String audioCachePath = "audioCache";

  // im文档路径
  static const String imDocumentPath = "imDocument";
}
