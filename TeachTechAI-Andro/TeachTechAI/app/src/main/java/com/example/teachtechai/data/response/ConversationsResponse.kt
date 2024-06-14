package com.example.teachtechai.data.response

import com.google.gson.annotations.SerializedName

data class TopicResponse(

	@field:SerializedName("data")
	val data: List<TopicItem?>? = null,

	@field:SerializedName("message")
	val message: String? = null,

	@field:SerializedName("errors")
	val errors: Any? = null,

	@field:SerializedName("status")
	val status: Boolean? = null
)

data class TopicItem(

	@field:SerializedName("start_time")
	val startTime: String? = null,

	@field:SerializedName("user_id")
	val userId: String? = null,

	@field:SerializedName("end_time")
	val endTime: String? = null,

	@field:SerializedName("topic")
	val topic: String? = null,

	@field:SerializedName("id")
	val id: String? = null
)
