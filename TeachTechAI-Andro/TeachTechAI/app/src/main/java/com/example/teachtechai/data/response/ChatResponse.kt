package com.example.teachtechai.data.response

import com.google.gson.annotations.SerializedName

data class ChatResponse(

	@field:SerializedName("data")
	val data: DataResponse? = null,

	@field:SerializedName("message")
	val message: String? = null,

	@field:SerializedName("errors")
	val errors: Any? = null,

	@field:SerializedName("status")
	val status: Boolean? = null,

	val isUser: Boolean
)

data class DataResponse(
	@field:SerializedName("conversation_id")
	val conversationId: String? = null,

	@field:SerializedName("response")
	val response: String? = null,

	@field:SerializedName("message_id")
	val messageId: String? = null
)

