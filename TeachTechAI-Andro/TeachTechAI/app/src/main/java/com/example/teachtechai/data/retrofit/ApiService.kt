package com.example.teachtechai.data.retrofit

import com.example.teachtechai.data.response.ChatResponse
import com.example.teachtechai.data.response.UpdateUserResponse
import com.example.teachtechai.data.response.GetMeResponse
import com.example.teachtechai.data.response.LoginResponse
import com.example.teachtechai.data.response.LogoutResponse
import com.example.teachtechai.data.response.OTPResponse
import com.example.teachtechai.data.response.RegisterResponse
import com.example.teachtechai.data.response.UploadProfileResponse
import com.example.teachtechai.data.response.VerifyOTPResponse
import okhttp3.MultipartBody
import retrofit2.Call
import retrofit2.http.Field
import retrofit2.http.FormUrlEncoded
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.Multipart
import retrofit2.http.PATCH
import retrofit2.http.POST
import retrofit2.http.Part

interface ApiService {
    @FormUrlEncoded
    @POST("user/register")
    fun register(
        @Field("name") name: String,
        @Field("email") email: String,
        @Field("password") password: String
    ): Call<RegisterResponse>
    @FormUrlEncoded
    @POST("user/login")
    fun login(
        @Field("email") email: String,
        @Field("password") password: String
    ): Call<LoginResponse>

    @POST("user/logout")
    fun logout(
        @Header("Authorization") token : String
    ): Call<LogoutResponse>

    @FormUrlEncoded
    @POST("user/forgot-password")
    fun forgetpassword(
        @Field("email") email : String
    ): Call<UpdateUserResponse>

    @FormUrlEncoded
    @POST("user/send-otp")
    fun sendotp(
        @Field("email") email : String
    ): Call<OTPResponse>

    @FormUrlEncoded
    @POST("user/verify-otp")
    fun verifyotp(
        @Field("email") email : String,
        @Field("otp") otp : String
    ): Call<VerifyOTPResponse>

    @FormUrlEncoded
    @PATCH("user/update")
    fun updateuser(
        @Header("Authorization") token : String,
        @Field("name") name : String,
        @Field("asal_instansi") asal_instansi : String,
        @Field("date_of_birth") date_of_birth : String
    ): Call<UpdateUserResponse>

    @GET("user/me")
    fun getme(
        @Header("Authorization") token : String
    ):Call<GetMeResponse>

    @FormUrlEncoded
    @POST("message/prompt")
    fun getChatResponse(
        @Header("Authorization") token : String,
        @Field("topic") topic: String,
        @Field("request") request: String,
        @Field("aimodel") aimodel: String,
    ): Call<ChatResponse>

    @Multipart
    @POST("user/upload-profile-picture")
    fun uploadProfile(
        @Header("Authorization") token : String,
        @Part photo : MultipartBody.Part
    ) : Call<UploadProfileResponse>
}