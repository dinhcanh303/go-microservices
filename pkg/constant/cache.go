package constant

import "time"

const (
	CachePrefix string = "mcs:social:"
)
const (
	CacheExpiresInOneSecond = time.Second
	CacheExpiresInOneMinute = CacheExpiresInOneSecond * 60
	CacheExpiresInOneHour   = CacheExpiresInOneMinute * 60
	CacheExpiresInOneDay    = CacheExpiresInOneHour * 24
	CacheExpiresInOneMonth  = CacheExpiresInOneDay * 30
	CacheExpiresInOneYear   = CacheExpiresInOneMonth * 12
)
const (
	CacheServiceAuth    = CachePrefix + "sv_auth:"
	CacheServicePost    = CachePrefix + "sv_post:"
	CacheServiceComment = CachePrefix + "sv_comment:"
	CacheServiceLike    = CachePrefix + "sv_like:"
	CacheServiceGroup   = CachePrefix + "sv_group:"
	CacheServiceUpload  = CachePrefix + "sv_upload:"
	CacheServiceSearch  = CachePrefix + "sv_search:"
)
const (
	CacheUsers        = CacheServiceAuth + "users"
	CacheKeyTokenUser = CacheServiceAuth + "key_token_user:"
)
const (
	CacheCommentsByCommentId   = CacheServiceComment + "comments:comment_id:"
	CacheCommentsByPostId      = CacheServiceComment + "comments:post_id:"
	CacheCommentsCountByPostId = CacheServiceComment + "comments:count:post_id:"
	CacheComments              = CacheServiceComment + "comments:"
)
const (
	CacheGroups           = CacheServiceGroup + "groups"
	CacheGroup            = CacheServiceGroup + "group:"
	CacheGroupsByUserId   = CacheServiceGroup + "groups:user_id:"
	CacheGroupIdsByUserId = CacheServiceGroup + "group_ids:user_id:"
	CacheGroupMembers     = CacheServiceGroup + "group_members:"
)
const (
	CacheLikeInfoByLikeableId = CacheServiceLike + "like_info:likeable_id:"
)
const (
	CachePosts          = CacheServicePost + "posts:"
	CachePostsFeed      = CacheServicePost + "posts:feed:"
	CachePostsFeedGroup = CacheServicePost + "posts:feed:group:"
	CachePostsGroupId   = CacheServicePost + "posts:group_id:"
	CachePostsUserId    = CacheServicePost + "posts:user_id:"
)
const (
	CacheAttachments = CacheServiceUpload + "attachments:"
)
const (
	CacheLimit  = ":limit:"
	CacheOffset = ":offset:"
	CacheUserId = ":user_id:"
)
