import 'dart:io';
import 'dart:math';
import 'package:client/common/common.dart';
import 'package:client/common/manage.dart';
import 'package:client/tool/file.dart';
import 'package:client/tool/log.dart';
import 'package:client/tool/tool.dart';
import 'package:minio/io.dart';
import 'package:minio/minio.dart';

/*
* oss管理
* @author wuxubaiyang
* @Time 2022/9/7 13:56
*/
class OSSManage extends BaseManage {
  static final OSSManage _instance = OSSManage._internal();

  factory OSSManage() => _instance;

  OSSManage._internal()
      : _minio = Minio(
          endPoint: OSSConfig.endPoint,
          port: OSSConfig.port,
          useSSL: OSSConfig.useSSL,
          accessKey: OSSConfig.accessKey,
          secretKey: OSSConfig.secretKey,
        );

  // oss单例
  final Minio _minio;

  // 上传附件
  Future<List<String?>> uploadFiles(
    List<File> files, {
    void Function(int)? onProgress,
    OSSBucket bucket = OSSBucket.jTechRecipe,
  }) async {
    final objects = <String?>[];
    for (final it in files) {
      try {
        final obj = _genObjectName(bucket.name, it);
        await _minio.fPutObject(bucket.name, obj, it.path);
        objects.add(obj);
        onProgress?.call(objects.length);
      } catch (e) {
        objects.add(null);
      }
    }
    return objects;
  }

  // 文件请求缓存
  final _objectGetCacheMap = <String, String>{};

  // 获取附件流
  Future<String> getObjectUrl(
    String object, {
    int? expires,
    OSSBucket bucket = OSSBucket.jTechRecipe,
    bool cached = true,
  }) async {
    try {
      if (object.isEmpty) return "";
      final cacheUrl = _objectGetCacheMap[object];
      if (cacheUrl != null) return cacheUrl;
      final url = await _minio.presignedGetObject(
        bucket.name,
        object,
        expires: expires,
      );
      if (cached) _objectGetCacheMap[object] = url;
      return url;
    } catch (e) {
      LogTool.e(e.toString());
    }
    return "";
  }

  // 生成附件对象名称
  String _genObjectName(String bucket, File file) {
    final name =
        '${file.path}_${Random(9527).nextDouble()}_${DateTime.now().toString()}';
    return '${bucket}_${Tool.md5(name)}${file.suffixes ?? ''}';
  }
}

// 单例调用
final ossManage = OSSManage();

// oss桶类型
enum OSSBucket { jTechRecipe }

// oss桶扩展
extension OSSBucketExtension on OSSBucket {
  // 获取桶的名称
  String get name => {
        OSSBucket.jTechRecipe: 'jtechrecipe',
      }[this]!;
}
