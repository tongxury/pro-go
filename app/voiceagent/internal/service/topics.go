package service

import (
	"context"
	pb "store/api/voiceagent"
)

// Static list of topics
var AllTopics = []*pb.Topic{
	{
		Id:          "anxiety",
		Title:       "缓解焦虑",
		Desc:        "平复心情，寻找内心的安宁",
		Greeting:    "你好，很高兴你愿意迈出这一步。刚才你提到有些缓解焦虑的需求，能和我聊聊最近发生了什么吗？",
		Instruction: "The user has selected to talk about anxiety. Focus your consultation specifically on anxiety relief. Use professional psychological techniques related to anxiety.",
	},
	{
		Id:          "stress",
		Title:       "压力管理",
		Desc:        "释放压力，重获掌控感",
		Greeting:    "你好，感觉到你有不少压力。想和我说说是工作还是生活上的事情吗？",
		Instruction: "The user has selected to talk about stress management. Focus on helping the user identify stressors and providing relief techniques.",
	},
	{
		Id:          "relationship",
		Title:       "人际关系",
		Desc:        "改善沟通，建立良性连接",
		Greeting:    "你好，关于人际关系的问题，通过倾诉往往能理清思路。你想聊聊具体的情况吗？",
		Instruction: "The user has selected to talk about relationships. Focus on interpersonal dynamics, communication, and social connection.",
	},
	{
		Id:          "mood",
		Title:       "情绪调节",
		Desc:        "接纳情绪，走出低谷",
		Greeting:    "你好，注意到你情绪有些低落。没关系，我们可以慢慢聊，发生了什么？",
		Instruction: "The user has selected to talk about low mood. Focus on empathy, validation, and exploring the causes of their mood.",
	},
	{
		Id:          "career",
		Title:       "职场困扰",
		Desc:        "突破瓶颈，明确职业方向",
		Greeting:    "你好，职场上的困扰确实挺消耗人的。你最近在工作中遇到了什么难题吗？",
		Instruction: "The user has selected to talk about career troubles. Focus on professional development, work-life balance, and career planning.",
	},
	{
		Id:          "intimate",
		Title:       "亲密关系",
		Desc:        "理解爱与被爱，化解矛盾",
		Greeting:    "你好，亲密关系的话题往往比较敏感，但也很重要。你愿意和我说说你的感受吗？",
		Instruction: "The user has selected to talk about intimate relations. Focus on love, marriage counseling, and deep emotional connection or conflict.",
	},
	{
		Id:          "growth",
		Title:       "自我成长",
		Desc:        "探索自我，实现潜能",
		Greeting:    "你好，很高兴看到你关注自我成长。最近有什么事情让你产生了这个念头吗？",
		Instruction: "The user has selected to talk about self-growth. Focus on exploring inner potential, self-awareness, and personal development goals.",
	},
	{
		Id:          "free",
		Title:       "自由对话",
		Desc:        "随心所欲，想聊什么都可以",
		Greeting:    "你好，我是 AURA。想聊点什么都行，我在这里陪你。",
		Instruction: "The user has selected free chat. Be open, supportive, and follow the user's lead while maintaining a professional counseling persona.",
	},
}

func (s *VoiceAgentService) ListTopics(ctx context.Context, req *pb.ListTopicsRequest) (*pb.TopicList, error) {
	return &pb.TopicList{
		List: AllTopics,
	}, nil
}
