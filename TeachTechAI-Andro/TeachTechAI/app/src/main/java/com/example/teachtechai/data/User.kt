package com.example.teachtechai.data

import com.bumptech.glide.load.model.GlideUrl

data class User(
    val id : String,
    val email : String,
    val name : String,
    val glideUrl : GlideUrl? = null,
    val nama_instansi : String? = null,
    val tanggal_lahir : String? = null,
)
