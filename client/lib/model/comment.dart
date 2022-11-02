import 'package:client/common/model.dart';
import 'package:client/model/model.dart';

/*
* 评论对象
* @author wuxubaiyang
* @Time 2022/10/14 15:05
*/
class CommentModel extends BaseModel with BasePart, CreatorPart {
  // 父id
  final String pId;

  // 评论类型
  final String typeCode;

  // 评论内容
  final String content;

  // 是否点赞
  final bool liked;

  // 点赞数量
  final num likeCount;

  CommentModel.from(obj)
      : pId = obj?['pId'] ?? '',
        typeCode = obj?['typeCode'] ?? '',
        content = obj?['content'] ?? '',
        liked = obj?['liked'] ?? false,
        likeCount = obj?['likeCount'] ?? 0 {
    initBasePart(obj);
    initCreatorPart(obj);
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        ...creatorPart,
        'pId': pId,
        'typeCode': typeCode,
        'content': content,
        'liked': liked,
        'likeCount': likeCount,
      };
}

/*
* 评论回复对象
* @author wuxubaiyang
* @Time 2022/10/14 15:05
*/
class ReplayModel extends BaseModel with BasePart, CreatorPart {
  // 父id
  final String pId;

  // 回复内容
  final String content;

  // 是否点赞
  final bool liked;

  // 点赞数量
  final num likeCount;

  ReplayModel.from(obj)
      : pId = obj?['pId'] ?? '',
        content = obj?['content'] ?? '',
        liked = obj?['liked'] ?? false,
        likeCount = obj?['likeCount'] ?? 0 {
    initBasePart(obj);
    initCreatorPart(obj);
  }

  @override
  Map<String, dynamic> to() => {
        ...basePart,
        ...creatorPart,
        'pId': pId,
        'content': content,
        'liked': liked,
        'likeCount': likeCount,
      };

  // 获取编辑结构
  Map<String, dynamic> toModifyInfo() => {
        'pId': pId,
        'content': content,
      };
}
