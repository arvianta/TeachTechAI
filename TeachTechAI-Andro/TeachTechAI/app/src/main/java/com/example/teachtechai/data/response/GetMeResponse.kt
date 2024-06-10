package com.example.teachtechai.data.response

import com.google.gson.annotations.SerializedName

data class GetMeResponse(

	@field:SerializedName("data")
	val data: Data? = null,

	@field:SerializedName("message")
	val message: String? = null,

	@field:SerializedName("errors")
	val errors: Any? = null,

	@field:SerializedName("status")
	val status: Boolean? = null
)

data class Data(

	@field:SerializedName("google_id")
	val googleId: String? = null,

	@field:SerializedName("date_of_birth")
	val dateOfBirth: String? = null,

	@field:SerializedName("created_at")
	val createdAt: String? = null,

	@field:SerializedName("profile_picture")
	val profilePicture: String? = null,

	@field:SerializedName("rt_expires")
	val rtExpires: String? = null,

	@field:SerializedName("session_token")
	val sessionToken: String? = null,

	@field:SerializedName("st_expires")
	val stExpires: String? = null,

	@field:SerializedName("refresh_token")
	val refreshToken: String? = null,

	@field:SerializedName("password")
	val password: String? = null,

	@field:SerializedName("updated_at")
	val updatedAt: String? = null,

	@field:SerializedName("phone")
	val phone: String? = null,

	@field:SerializedName("role_id")
	val roleId: String? = null,

	@field:SerializedName("name")
	val name: String? = null,

	@field:SerializedName("asal_instansi")
	val asalInstansi: String? = null,

	@field:SerializedName("id")
	val id: String? = null,

	@field:SerializedName("DeletedAt")
	val deletedAt: Any? = null,

	@field:SerializedName("email")
	val email: String? = null
)
