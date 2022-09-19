import 'package:client/common/model.dart';
import 'package:client/model/user.dart';

import 'model.dart';
import 'post_content.dart';

/*
* 帖子实体
* @author wuxubaiyang
* @Time 2022/9/12 21:49
*/
class PostModel extends BaseModel with BaseInfo, CreatorInfo {
  // 标题
  final String title;

  // 内容集合
  final List<PostContentItem> contents;

  // 浏览数
  final num viewCount;

  // 是否浏览
  final bool viewed;

  // 点赞人集合
  final List<UserModel> likeUserList;

  // 点赞数
  final num likeCount;

  // 是否点赞
  final bool liked;

  // 收藏数
  final num collectCount;

  // 是否收藏
  final bool collected;

  // 标签集合
  final List<PostTagModel> tags;

  PostModel.from(obj)
      : title = obj?["title"] ?? "",
        contents = (obj?["contents"] ?? [])
            .map<PostContentItem>((e) => PostContentItem.from(e))
            .toList(),
        viewCount = obj?["viewCount"] ?? 0,
        viewed = obj?["viewed"] ?? false,
        likeUserList = (obj?["likeUserList"] ?? [])
            .map<UserModel>((e) => UserModel.from(e))
            .toList(),
        likeCount = obj?["likeCount"] ?? 0,
        liked = obj?["liked"] ?? false,
        collectCount = obj?["collectCount"] ?? 0,
        collected = obj?["collected"] ?? false,
        tags = (obj?["tags"] ?? [])
            .map<PostTagModel>((e) => PostTagModel.from(e))
            .toList() {
    initialBaseInfo(obj);
    initialCreatorInfo(obj);
  }

  @override
  Map<String, dynamic> to() => {
        "title": title,
        "contents": contents.map((e) => e.to()).toList(),
        "viewCount": viewCount,
        "viewed": viewed,
        "likeUserList": likeUserList.map((e) => e.to()).toList(),
        "likeCount": likeCount,
        "liked": liked,
        "collectCount": collectCount,
        "collected": collected,
        "tags": tags.map((e) => e.to()).toList(),
      };
}

/*
* 帖子评论
* @author wuxubaiyang
* @Time 2022/9/12 20:56
*/
class PostCommentModel extends BaseModel with BaseInfo, CreatorInfo {
  // 帖子id
  final num postId;

  // 评论内容
  final String content;

  // 回复用户集合
  final List<UserModel> replays;

  // 回复数量
  final num replayCount;

  // 点赞数量
  final num likeCount;

  // 是否点赞
  final bool liked;

  PostCommentModel.from(obj)
      : postId = obj?["postId"] ?? 0,
        content = obj?["content"] ?? "",
        replays = (obj?["replays"] ?? [])
            .map<PostCommentReplayModel>((e) => PostCommentReplayModel.from(e))
            .toList(),
        replayCount = obj?["replayCount"] ?? 0,
        likeCount = obj?["likeCount"] ?? 0,
        liked = obj?["liked"] ?? false {
    initialBaseInfo(obj);
    initialCreatorInfo(obj);
  }
}

/*
* 帖子评论回复
* @author wuxubaiyang
* @Time 2022/9/12 20:56
*/
class PostCommentReplayModel extends BaseModel with BaseInfo, CreatorInfo {
  // 评论id
  final num commentId;

  // 回复内容
  final String content;

  // 点赞数量
  final num likeCount;

  // 是否点赞
  final bool liked;

  PostCommentReplayModel.from(obj)
      : commentId = obj?["commentId"] ?? 0,
        content = obj?["content"] ?? "",
        likeCount = obj?["likeCount"] ?? 0,
        liked = obj?["liked"] ?? false {
    initialBaseInfo(obj);
    initialCreatorInfo(obj);
  }
}

/*
* 帖子标签实体
* @author wuxubaiyang
* @Time 2022/9/12 20:50
*/
class PostTagModel extends BaseModel with BaseInfo, CreatorInfo {
  // 帖子id
  final num postId;

  // 标签名称
  final String name;

  PostTagModel.from(obj)
      : postId = obj?["postId"] ?? 0,
        name = obj?["name"] ?? "" {
    initialBaseInfo(obj);
    initialCreatorInfo(obj);
  }

  @override
  Map<String, dynamic> to() => {
        ...baseInfoMap,
        ...creatorInfoMap,
        "postId": postId,
        "name": name,
      };
}
