package com.example.teachtechai.di

import android.content.Context
import com.example.teachtechai.data.UserRepository
import com.example.teachtechai.data.pref.UserPreference
import com.example.teachtechai.data.pref.dataStore


object Injection {
    fun provideRepository(context: Context): UserRepository {
        val pref = UserPreference.getInstance(context.dataStore)
        return UserRepository.getInstance(pref)
    }
}