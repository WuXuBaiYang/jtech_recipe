import 'dart:io';
import 'dart:math';
import 'package:client/common/common.dart';
import 'package:client/common/manage.dart';
import 'package:client/tool/file.dart';
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
    required String bucket,
    void Function(int)? onProgress,
  }) async {
    var objects = <String?>[];
    for (var it in files) {
      try {
        var obj = _genObjectName(bucket, it);
        await _minio.fPutObject(bucket, obj, it.path);
        objects.add(obj);
        onProgress?.call(objects.length);
      } catch (e) {
        objects.add(null);
      }
    }
    return objects;
  }

  // 上传附件到Profile
  Future<List<String?>> uploadProfileFile(
    List<File> files, {
    void Function(int)? onProgress,
  }) =>
      uploadFiles(files, bucket: "profile", onProgress: onProgress);

  // 上传附件到Post
  Future<List<String?>> uploadPostFile(
    List<File> files, {
    void Function(int)? onProgress,
  }) =>
      uploadFiles(files, bucket: "post", onProgress: onProgress);

  // 获取附件流
  Future<String> getObjectUrl(
    String object, {
    required String bucket,
    int? expires,
  }) =>
      _minio.presignedGetObject(bucket, object, expires: expires);

  // 获取Profile附件流
  Future<String> getProfileObjectUrl(
    String object, {
    int? expires,
  }) =>
      getObjectUrl(object, bucket: "profile", expires: expires);

  // 获取Post附件流
  Future<String> getPostObjectUrl(
    String object, {
    int? expires,
  }) =>
      getObjectUrl(object, bucket: "post", expires: expires);

  // 生成附件对象名称
  String _genObjectName(String bucket, File file) {
    var name = "${file.path}_${Random(9527).nextDouble()}";
    return "${bucket}_${Tool.md5(name)}${file.suffixes ?? ""}";
  }
}

// 单例调用
final ossManage = OSSManage();
