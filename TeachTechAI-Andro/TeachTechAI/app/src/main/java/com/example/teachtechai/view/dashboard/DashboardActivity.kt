package com.example.teachtechai.view.dashboard

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import androidx.navigation.NavController
import androidx.navigation.Navigation
import androidx.navigation.fragment.NavHostFragment
import com.example.teachtechai.MainActivity
import com.example.teachtechai.R
import com.example.teachtechai.data.pref.UserPreference
import com.example.teachtechai.data.pref.dataStore
import com.example.teachtechai.databinding.ActivityDashboardBinding
import com.google.android.material.bottomnavigation.BottomNavigationView
import kotlinx.coroutines.runBlocking

class DashboardActivity : AppCompatActivity() {
    private lateinit var binding: ActivityDashboardBinding
    private lateinit var navController: NavController
    private lateinit var userPreference : UserPreference
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityDashboardBinding.inflate(layoutInflater)
        setContentView(binding.root)

        userPreference = UserPreference.getInstance(this.dataStore)
        redirectToLogin()
        navController = Navigation.findNavController(this, R.id.dashboard_container)
        val bottomNavView = findViewById<BottomNavigationView>(R.id.nav_view)
        bottomNavView.setOnNavigationItemSelectedListener { item ->
            when(item.itemId){
                R.id.navigation_discover -> {
                    navController.navigate(R.id.discoverFragment)
                    true
                }
                R.id.navigation_search ->{
                    navController.navigate(R.id.searchFragment)
                    true
                }
                R.id.navigation_profile ->{
                    navController.navigate(R.id.profileFragment)
                    true
                }
                else -> false
            }
        }
    }
    private fun redirectToLogin() {
        runBlocking {
            val token = userPreference.getToken()
            if(token == null){
                navigateToLogin()
            }
        }
    }
    private fun navigateToLogin(){
        val intent = Intent(this, MainActivity::class.java)
        startActivity(intent)
    }

}