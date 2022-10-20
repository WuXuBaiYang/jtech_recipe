import 'package:client/api/comment.dart';
import 'package:client/api/tag.dart';
import 'package:client/common/api/request.dart';
import 'package:client/model/comment.dart';
import 'package:client/model/model.dart';
import 'package:client/model/recipe.dart';
import 'package:client/model/tag.dart';
import 'base.dart';

/*
* 食谱接口
* @author wuxubaiyang
* @Time 2022/9/12 18:48
*/
class RecipeAPI extends BaseJAPI {
  // 创建食谱
  Future<RecipeModel> createRecipe({
    required RecipeModel model,
  }) {
    return handleResponseData(
      post("/recipe",
          requestModel: RequestModel.body(
            data: model.toModifyInfo(),
          )),
      handle: (e) => RecipeModel.from(e),
    );
  }

  // 编辑食谱
  Future<RecipeModel> updateRecipe({
    required RecipeModel model,
    required String recipeId,
  }) {
    return handleResponseData(
      put("/recipe/$recipeId",
          requestModel: RequestModel.body(
            data: model.toModifyInfo(),
          )),
      handle: (e) => RecipeModel.from(e),
    );
  }

  // 获取食谱集合
  Future<PaginationModel<RecipeModel>> loadMenus({
    int pageIndex = 1,
    int pageSize = 15,
    String? userId,
  }) {
    return handleResponsePaginationData(
      get("/recipe",
          requestModel: RequestModel.query(
            parameters: {
              "pageIndex": pageIndex,
              "pageSize": pageSize,
              if (userId != null) "userId": userId,
            },
          )),
      handle: (e) => RecipeModel.from(e),
    );
  }

  // 获取食谱详情
  Future<RecipeModel> loadRecipeInfo({
    required String recipeId,
  }) {
    return handleResponseData(
      get("/recipe/$recipeId"),
      handle: (e) => RecipeModel.from(e),
    );
  }

  // 发布食谱评论
  Future<CommentModel> createRecipeComment({
    required String recipeId,
    required String content,
  }) {
    return commentApi.createComment(
      path: "/recipe/comment",
      pId: recipeId,
      content: content,
    );
  }

  // 获取食谱评论列表
  Future<PaginationModel<CommentModel>> loadRecipeComments({
    required String recipeId,
    int pageIndex = 1,
    int pageSize = 15,
  }) {
    return commentApi.loadComments(
      path: "/recipe/comment",
      pId: recipeId,
      pageIndex: pageIndex,
      pageSize: pageSize,
    );
  }

  // 批量添加食谱标签
  Future<List<TagModel>> addRecipeTags({
    required List<TagModel> tags,
  }) {
    return tagApi.addTags(
      path: "/recipe/tag",
      tags: tags,
    );
  }

  // 获取食谱标签集合
  Future<PaginationModel<TagModel>> loadRecipeTags({
    int pageIndex = 1,
    int pageSize = 15,
    String? userId,
  }) {
    return tagApi.loadTags(
      path: "/recipe/tag",
      pageIndex: pageIndex,
      pageSize: pageSize,
      userId: userId,
    );
  }

  // 食谱点赞
  Future<bool> likeRecipe({
    required String recipeId,
  }) {
    return handleResponseData(
      post("/recipe/like/$recipeId"),
    );
  }

  // 食谱取消点赞
  Future<bool> unLikeRecipe({
    required String recipeId,
  }) {
    return handleResponseData(
      delete("/recipe/like/$recipeId"),
    );
  }

  // 食谱收藏
  Future<bool> collectRecipe({
    required String recipeId,
  }) {
    return handleResponseData(
      post("/recipe/collect/$recipeId"),
    );
  }

  // 食谱取消收藏
  Future<bool> unCollectRecipe({
    required String recipeId,
  }) {
    return handleResponseData(
      delete("/recipe/collect/$recipeId"),
    );
  }
}

// 单例调用
final recipeApi = RecipeAPI();
