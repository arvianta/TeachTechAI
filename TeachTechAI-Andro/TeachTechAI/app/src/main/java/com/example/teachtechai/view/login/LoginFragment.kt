package com.example.teachtechai.view.login

import android.content.Intent
import android.os.Bundle
import android.text.Editable
import android.text.InputType
import android.text.TextWatcher
import android.util.Log
import android.util.Patterns
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.EditText
import android.widget.ImageView
import androidx.activity.OnBackPressedCallback
import androidx.activity.addCallback
import androidx.fragment.app.viewModels
import androidx.lifecycle.lifecycleScope
import com.example.teachtechai.R
import com.example.teachtechai.data.pref.UserPreference
import com.example.teachtechai.data.pref.dataStore
import com.example.teachtechai.databinding.FragmentLoginBinding
import com.example.teachtechai.view.ViewModelFactory
import com.example.teachtechai.view.dashboard.DashboardActivity
import com.example.teachtechai.view.forgetpassword.ForgetPassword
import com.example.teachtechai.view.register.RegisterFragment
import kotlinx.coroutines.launch

class LoginFragment : Fragment() {
    private lateinit var binding: FragmentLoginBinding
    private var isVisiblity = false
    private val loginViewModel by viewModels<LoginViewModel> {
        ViewModelFactory.getInstance(requireContext())
    }
    private lateinit var userPreference: UserPreference

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        binding = FragmentLoginBinding.inflate(inflater, container,false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        userPreference = UserPreference.getInstance(requireContext().dataStore)

        emailpasswordValidation()
        navigatetoLupaKataSandi()
        navigateToDaftarSekarang()

        buttonLogin()
        observeData()

        requireActivity().onBackPressedDispatcher.addCallback(viewLifecycleOwner, object : OnBackPressedCallback(true) {
            override fun handleOnBackPressed() {
                requireActivity().finish()
            }
        })
    }

    private fun observeData() {
        loginViewModel.loginResult.observe(viewLifecycleOwner) { response ->
            if (response.status == true) {
                navigateToDashboard()
            }
        }
        loginViewModel.isLoading.observe(viewLifecycleOwner) { isLoading ->
            binding.loginProgressBar.visibility = if (isLoading) View.VISIBLE else View.GONE
        }

        loginViewModel.errorMessage.observe(viewLifecycleOwner){errorMessage ->
            Log.d("ERROR MESSAGE", "$errorMessage")
            if (errorMessage != null) {
                binding.loginTvKataSandiSalah.visibility = View.VISIBLE
            }
        }
    }

    private fun emailpasswordValidation() {
        val emailEditText = binding.loginEditEmail
        val passwordEditText = binding.loginEditPassword
        val passwordToggle = binding.passwordToggle

        emailEditText.addTextChangedListener(object : TextWatcher {
            override fun afterTextChanged(s: Editable?) {
                validateEmail(emailEditText)
            }

            override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) {}
            override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) {}
        })
        passwordEditText.addTextChangedListener(object : TextWatcher {
            override fun afterTextChanged(s: Editable?) {
                validatePassword(passwordEditText)
            }

            override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) {}
            override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) {}
        })
        passwordToggle.setOnClickListener{
            togglePasswordVisiblity(passwordEditText, passwordToggle)
        }
    }

    private fun checkToken(){
       lifecycleScope.launch {
            val token = userPreference.getToken()
            if(token != null){
                navigateToDashboard()
            }
       }
    }

    private fun buttonLogin(){
        binding.loginButtonMasuk.setOnClickListener {
            val email = binding.loginEditEmail.text.toString()
            val password = binding.loginEditPassword.text.toString()

            if(email.isNotEmpty() && password.isNotEmpty()){
                loginViewModel.loginUser(email, password)
            }
        }

    }
    private fun navigateToDashboard() {
        val intent = Intent(requireContext(), DashboardActivity :: class.java)
        startActivity(intent)
    }

    private fun navigatetoLupaKataSandi() {
        binding.loginLupaKataSandi.setOnClickListener{
            parentFragmentManager.beginTransaction()
                .replace(R.id.fragment_container, ForgetPassword())
                .addToBackStack(null)
                .commit()
        }
    }
    private fun navigateToDaftarSekarang() {
        binding.loginTvDaftarSekarang.setOnClickListener {
            parentFragmentManager.beginTransaction()
                .replace(R.id.fragment_container, RegisterFragment())
                .addToBackStack(null)
                .commit()
        }
    }

    private fun validateEmail(emailEditText : EditText) {
        val email = emailEditText.text.toString()
        if(email.isEmpty()){
            emailEditText.error = null
        }else if (!Patterns.EMAIL_ADDRESS.matcher(email).matches()) {
            emailEditText.error = "Email tidak valid"
        } else {
            emailEditText.error = null
        }
    }
    private fun validatePassword(passwordEditText : EditText) {
        val password = passwordEditText.text.toString()
        if(password.isEmpty()){
            passwordEditText.error = null
        }else if (password.length < 8) {
            passwordEditText.error = "Kata sandi minimal 8 karakter"
        } else {
            passwordEditText.error = null
        }
    }
    private fun togglePasswordVisiblity(passwordEditText: EditText, passwordToggle : ImageView){
        if(isVisiblity == true){
            passwordEditText.inputType = InputType.TYPE_TEXT_VARIATION_PASSWORD or InputType.TYPE_CLASS_TEXT
            passwordToggle.setImageResource(R.drawable.baseline_visibility_24)
        }else{
            passwordEditText.inputType = InputType.TYPE_TEXT_VARIATION_VISIBLE_PASSWORD
            passwordToggle.setImageResource(R.drawable.baseline_visibility_off_24)
        }
        passwordEditText.setSelection(passwordEditText.text.length)
        isVisiblity = !isVisiblity
    }
}