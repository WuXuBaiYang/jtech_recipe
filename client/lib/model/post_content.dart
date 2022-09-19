import 'package:client/common/model.dart';

/*
* 帖子内容项结构体
* @author wuxubaiyang
* @Time 2022/9/12 21:06
*/
abstract class PostContentItem extends BaseModel {
  // 内容类型
  final PostContentType type;

  // 创建文本实体
  static PostContentTextItem text(String content) =>
      PostContentTextItem.create(content);

  // 创建图片实体
  static PostContentImageItem image(List<String> images) =>
      PostContentImageItem.create(images);

  // 创建视频实体
  static PostContentVideoItem video(String videoUrl) =>
      PostContentVideoItem.create(videoUrl);

  // 创建音频实体
  static PostContentAudioItem audio(String audioUrl) =>
      PostContentAudioItem.create(audioUrl);

  // 从json中解析
  static PostContentItem from(obj) {
    var type = PostContentType.values[obj?["type"] ?? 0];
    switch (type) {
      case PostContentType.text:
        return PostContentTextItem.from(obj);
      case PostContentType.image:
        return PostContentImageItem.from(obj);
      case PostContentType.video:
        return PostContentVideoItem.from(obj);
      case PostContentType.audio:
        return PostContentAudioItem.from(obj);
    }
  }

  // 主构造
  PostContentItem(this.type);

  @override
  Map<String, dynamic> to() => {
        "type": type.index,
      };
}

/*
* 帖子内容项-文本
* @author wuxubaiyang
* @Time 2022/9/12 21:14
*/
class PostContentTextItem extends PostContentItem {
  // 文本内容
  final String content;

  PostContentTextItem.create(this.content) : super(PostContentType.text);

  PostContentTextItem.from(obj)
      : content = obj?["content"] ?? "",
        super(PostContentType.text);

  @override
  Map<String, dynamic> to() => {
        ...super.to(),
        "content": content,
      };
}

/*
* 帖子内容项-图片
* @author wuxubaiyang
* @Time 2022/9/12 21:14
*/
class PostContentImageItem extends PostContentItem {
  // 图片集合
  final List<String> images;

  PostContentImageItem.create(this.images) : super(PostContentType.image);

  PostContentImageItem.from(obj)
      : images = (obj?["images"] ?? []).map<String>((e) => e).toList(),
        super(PostContentType.image);

  @override
  Map<String, dynamic> to() => {
        ...super.to(),
        "images": images,
      };
}

/*
* 帖子内容项-视频
* @author wuxubaiyang
* @Time 2022/9/12 21:39
*/
class PostContentVideoItem extends PostContentItem {
  // 视频地址
  final String videoUrl;

  PostContentVideoItem.create(this.videoUrl) : super(PostContentType.video);

  PostContentVideoItem.from(obj)
      : videoUrl = obj?["videoUrl"] ?? "",
        super(PostContentType.video);

  @override
  Map<String, dynamic> to() => {
        ...super.to(),
        "videoUrl": videoUrl,
      };
}

/*
* 帖子内容项-音频
* @author wuxubaiyang
* @Time 2022/9/12 21:39
*/
class PostContentAudioItem extends PostContentItem {
  // 视频地址
  final String audioUrl;

  PostContentAudioItem.create(this.audioUrl) : super(PostContentType.audio);

  PostContentAudioItem.from(obj)
      : audioUrl = obj?["audioUrl"] ?? "",
        super(PostContentType.audio);

  @override
  Map<String, dynamic> to() => {
        ...super.to(),
        "audioUrl": audioUrl,
      };
}

/*
* 帖子内容项枚举
* @author wuxubaiyang
* @Time 2022/9/12 21:08
*/
enum PostContentType { text, image, video, audio }

/*
* 帖子内容项枚举扩展
* @author wuxubaiyang
* @Time 2022/9/12 21:09
*/
extension PostContentTypeExtension on PostContentType {
// 获取类型名称
  String get name {
    switch (this) {
      case PostContentType.text:
        return "文本";
      case PostContentType.image:
        return "图片";
      case PostContentType.video:
        return "视频";
      case PostContentType.audio:
        return "音频";
    }
  }
}
