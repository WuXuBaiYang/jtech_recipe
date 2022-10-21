import 'package:client/main.dart';
import 'package:client/page/home.dart';
import 'package:client/page/auth.dart';
import 'package:flutter/widgets.dart';

/*
* 路由路径静态变量
* @author wuxubaiyang
* @Time 2022/9/8 14:55
*/
class RoutePath {
  // 创建路由表
  static Map<String, WidgetBuilder> get routes => {
        home: (c) => const HomePage(),
        auth: (c) => const AuthPage(),
      };

  // 首页
  static const String home = '/home';

  // 授权页
  static const String auth = '/auth';
}

/*
* 静态资源/通用静态变量
* @author wuxubaiyang
* @Time 2022/9/8 14:54
*/
class Common {
  // api基础地址
  static String baseUrl = debugMode ? _baseUrlDev : _baseUrlRelease;

  // api开发地址
  static const String _baseUrlDev = 'http://192.168.16.50:9527/api';

  // api正式地址
  static const String _baseUrlRelease = 'xxxxx';

  // 加密盐
  static const String salt = 'p5bDTyO6McoYeH';
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
  static const String _endPointDev = '192.168.16.50';

  // 正式版地址
  static const String _endPointRelease = '';

  // 端口号
  static const int port = 9000;

  // 是否启用ssl
  static const useSSL = false;

  // 授权key
  static const accessKey = r'IkcLOGuLGKZRIN7b';

  // 密钥
  static const secretKey = r'VvZPcT8SOPl5pQDDshMqLGvpbvIeIhpo';
}

/*
* 目录管理
* @author wuxubaiyang
* @Time 2022/9/9 17:57
*/
class FileDirPath {
  // 图片缓存路径
  static const String imageCachePath = 'imageCache';

  // 视频缓存路径
  static const String videoCachePath = 'videoCache';

  // 音频缓存路径
  static const String audioCachePath = 'audioCache';

  // im文档路径
  static const String imDocumentPath = 'imDocument';
}
