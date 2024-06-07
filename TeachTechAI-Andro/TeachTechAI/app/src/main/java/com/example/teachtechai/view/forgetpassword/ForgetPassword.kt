package com.example.teachtechai.view.forgetpassword

import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import com.example.teachtechai.R
import com.example.teachtechai.databinding.FragmentForgetPasswordBinding
import com.example.teachtechai.view.inputotp.VerifyOtp

class ForgetPassword : Fragment() {
    private lateinit var binding : FragmentForgetPasswordBinding
    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        binding = FragmentForgetPasswordBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        binding.fpButtonReset.setOnClickListener {
            parentFragmentManager.beginTransaction()
                .replace(R.id.fragment_container, VerifyOtp())
                .addToBackStack(null)
                .commit()
        }
    }
}