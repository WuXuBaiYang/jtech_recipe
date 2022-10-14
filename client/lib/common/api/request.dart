import 'package:client/common/model.dart';
import 'package:dio/dio.dart';
import 'package:http_parser/http_parser.dart';

/*
* 请求对象实体
* @author wuxubaiyang
* @Time 2022/3/29 14:20
*/
class RequestModel extends BaseModel {
  // 查询参数
  final Map<String, dynamic>? parameters;

  // 头部参数
  final Map<String, dynamic>? headers;

  // 消息体
  final dynamic data;

  const RequestModel({
    this.parameters,
    this.headers,
    this.data,
  });

  // 创建构造
  const RequestModel.create({
    Map<String, dynamic>? parameters,
    Map<String, dynamic>? headers,
    dynamic data,
  }) : this(
          parameters: parameters,
          headers: headers,
          data: data,
        );

  // 构造为查询格式
  const RequestModel.query({
    required Map<String, dynamic>? parameters,
    Map<String, dynamic>? headers,
  }) : this(
          parameters: parameters,
          headers: headers,
          data: null,
        );

  // 构造为消息体格式
  const RequestModel.body({
    required dynamic data,
    Map<String, dynamic>? parameters,
    Map<String, dynamic>? headers,
  }) : this(
          parameters: parameters,
          headers: headers,
          data: data,
        );

  // 表单构建模式
  static RequestFormBuilder form({
    Map<String, dynamic>? parameters,
    Map<String, dynamic>? headers,
    Map<String, dynamic>? data,
  }) =>
      RequestFormBuilder(
        parameters: parameters,
        headers: headers,
        data: data,
      );
}

/*
* 请求实体表单构造对象
* @author wuxubaiyang
* @Time 2022/3/29 14:25
*/
class RequestFormBuilder extends RequestModel {
  RequestFormBuilder({
    Map<String, dynamic>? parameters,
    Map<String, dynamic>? headers,
    Map<String, dynamic>? data,
  }) : super(
            parameters: parameters,
            headers: headers,
            data: FormData.fromMap(data ?? {}));

  // 添加参数
  RequestFormBuilder add(String key, dynamic value) =>
      this..addAll({key: value});

  // 添加多个参数
  RequestFormBuilder addAll(Map<String, dynamic> data) => this
    ..data.fields.addAll(
        data.map((key, value) => MapEntry(key, value.toString())).entries);

  // 添加文件
  RequestFormBuilder addFileSync(
    String key,
    String filePath, {
    String? filename,
    MediaType? mediaType,
  }) {
    return this
      ..data.files.add(MapEntry(
            key,
            MultipartFile.fromFileSync(
              filePath,
              filename: filename,
              contentType: mediaType,
            ),
          ));
  }

  // 添加多个文件
  RequestFormBuilder addFilesSync(
    String key,
    List<RequestFileItem> files,
  ) =>
      this
        ..data.files.addAll(files
            .map((item) => MapEntry(
                key,
                MultipartFile.fromFileSync(
                  item.filePath,
                  filename: item.filename,
                  contentType: item.mediaType,
                )))
            .toList());

  // 构建为请求对象
  RequestModel build() => this;
}

/*
* 表单附件对象
* @author wuxubaiyang
* @Time 2022/3/29 14:29
*/
class RequestFileItem {
  // 文件路径
  final String filePath;

  // 文件名
  final String? filename;

  // 文件类型
  final MediaType? mediaType;

  RequestFileItem({
    required this.filePath,
    this.filename,
    this.mediaType,
  });
}
