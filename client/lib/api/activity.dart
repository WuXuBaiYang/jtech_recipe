import 'package:client/api/comment.dart';
import 'package:client/model/activity.dart';
import 'package:client/model/comment.dart';
import 'package:client/model/model.dart';
import 'base.dart';

/*
* 活动接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class ActivityAPI extends BaseJAPI {
  // 获取全部进行中的活动
  Future<List<ActivityRecordModel>> loadAllProcessActivity() {
    return handleResponseListData(
      get('/activity/process'),
      handle: (e) => ActivityRecordModel.from(e),
    );
  }

  // 发布活动评论
  Future<CommentModel> createActivityComment({
    required String activityId,
    required String content,
  }) {
    return commentApi.createComment(
      path: '/activity/comment',
      pId: activityId,
      content: content,
    );
  }

  // 获取活动评论列表
  Future<PaginationModel<CommentModel>> loadActivityComments({
    required String activityId,
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return commentApi.loadComments(
      path: '/activity/comment',
      pId: activityId,
      pageIndex: pageIndex,
      pageSize: pageSize,
    );
  }
}

// 单例调用
final activityApi = ActivityAPI();
