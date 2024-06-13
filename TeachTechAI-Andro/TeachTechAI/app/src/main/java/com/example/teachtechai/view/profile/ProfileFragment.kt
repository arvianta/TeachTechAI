package com.example.teachtechai.view.profile

import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.util.TypedValue
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Button
import androidx.appcompat.app.AlertDialog
import androidx.fragment.app.activityViewModels
import androidx.fragment.app.viewModels
import androidx.navigation.fragment.findNavController
import com.example.teachtechai.MainActivity
import com.example.teachtechai.R
import com.example.teachtechai.data.User
import com.example.teachtechai.data.pref.UserPreference
import com.example.teachtechai.data.pref.dataStore
import com.example.teachtechai.databinding.FragmentProfileBinding
import com.example.teachtechai.view.SharedViewModel
import com.example.teachtechai.view.ViewModelFactory
import com.example.teachtechai.view.discover.DiscoverViewModel
import com.example.teachtechai.view.editprofile.EditProfileFragment
import kotlinx.coroutines.runBlocking

class ProfileFragment : Fragment() {
    private lateinit var userPreference: UserPreference
    private lateinit var binding: FragmentProfileBinding
    private val logoutViewModel by viewModels<LogoutViewModel> {
        ViewModelFactory.getInstance(requireContext())
    }
    private val sharedViewModel : SharedViewModel by activityViewModels()
    private var alertDialog: AlertDialog? = null

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        // Inflate the layout for this fragment
        binding = FragmentProfileBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)
        userPreference = UserPreference.getInstance(requireContext().dataStore)
        checkToken()
        observeData()
        binding.profileTvKeluar.setOnClickListener {
            showDialogBoxRegister()
        }
        navigateToEditProfile()
    }

    private fun navigateToEditProfile() {
        binding.profileEdit.setOnClickListener {
            findNavController().navigate(R.id.profileFragment_to_editProfileFragment)
        }
    }

    private fun logout(){
        runBlocking {
            val token = userPreference.getToken()
            if (token != null) {
                logoutViewModel.logoutUser(token)
            }
        }
    }
    private fun showDialogBoxRegister(){
        val dialogView = LayoutInflater.from(requireContext()).inflate(R.layout.dialog_logout, null)
        val dialogBuilder = AlertDialog.Builder(requireContext())
            .setView(dialogView)

        val alertDialog = dialogBuilder.create()
        alertDialog.window?.setBackgroundDrawableResource(android.R.color.transparent)
        alertDialog.show()

        val width = TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, 301f, resources.displayMetrics).toInt()
        val height = TypedValue.applyDimension(TypedValue.COMPLEX_UNIT_DIP, 395f, resources.displayMetrics).toInt()
        alertDialog.window?.setLayout(width, height)
        val buttonKeluar = dialogView.findViewById<Button>(R.id.buttonKeluar)
        val buttonBatal = dialogView.findViewById<Button>(R.id.buttonBatal)
        buttonBatal.setOnClickListener {
            alertDialog.dismiss()
        }
        buttonKeluar.setOnClickListener {
            logout()
            alertDialog.dismiss()
            navigateToMainActivity()
        }
    }
    private fun checkToken(){
        runBlocking {
            val token = userPreference.getToken()
            if (token == null){
                navigateToMainActivity()
            }
        }
    }

    private fun observeData(){
        Log.d("OBSERVE DATA", "MASUK KESINI")
        sharedViewModel.user.observe(viewLifecycleOwner){user->
            Log.d("NAME PROFILE", "$user.name")
            binding.profileTvName.text = user.name
        }
    }
    private fun navigateToMainActivity(){
        val intent = Intent(requireContext(), MainActivity::class.java)
        intent.flags = Intent.FLAG_ACTIVITY_CLEAR_TASK or Intent.FLAG_ACTIVITY_NEW_TASK
        startActivity(intent)
    }

    override fun onDestroy() {
        super.onDestroy()
        alertDialog?.dismiss()
        alertDialog = null
    }

}