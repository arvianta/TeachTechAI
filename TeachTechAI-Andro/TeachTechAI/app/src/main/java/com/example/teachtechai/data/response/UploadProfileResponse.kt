package com.example.teachtechai.data.response

import com.google.gson.annotations.SerializedName

data class UploadProfileResponse(
	@field:SerializedName("data")
	val data: Data? = null,

	@field:SerializedName("message")
	val message: String? = null,

	@field:SerializedName("errors")
	val errors: String? = null,

	@field:SerializedName("status")
	val status: Boolean? = null
)

data class ChangePasswordResponse(
	@field:SerializedName("data")
	val data: Data? = null,

	@field:SerializedName("message")
	val message: String? = null,

	@field:SerializedName("errors")
	val errors: String? = null,

	@field:SerializedName("status")
	val status: Boolean? = null
)