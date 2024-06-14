package com.example.teachtechai.view.discover

import android.content.Intent
import android.os.Bundle
import android.util.Log
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.fragment.app.activityViewModels
import androidx.fragment.app.viewModels
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.GridLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import com.bumptech.glide.load.model.GlideUrl
import com.bumptech.glide.load.model.LazyHeaders
import com.bumptech.glide.load.resource.bitmap.CenterCrop
import com.bumptech.glide.load.resource.bitmap.CircleCrop
import com.example.teachtechai.MainActivity
import com.example.teachtechai.R
import com.example.teachtechai.data.User
import com.example.teachtechai.data.pref.UserPreference
import com.example.teachtechai.data.pref.dataStore
import com.example.teachtechai.databinding.FragmentDiscoverBinding
import com.example.teachtechai.databinding.FragmentDiscoverShimmerBinding
import com.example.teachtechai.view.SharedViewModel
import com.facebook.shimmer.ShimmerFrameLayout
import kotlinx.coroutines.launch
import kotlinx.coroutines.runBlocking

class DiscoverFragment : Fragment() {
    private lateinit var binding: FragmentDiscoverBinding
    private lateinit var userPreference : UserPreference
    private val discoverViewModel : DiscoverViewModel by viewModels()
    private val sharedViewModel : SharedViewModel by activityViewModels()
    var glideUrl : GlideUrl? = null

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        binding = FragmentDiscoverBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        userPreference = UserPreference.getInstance(requireContext().dataStore)
        val topics = listOf(
            Topic("SENI DAN OLAHRAGA", "Mengembangkan Kreativitas melalui Seni", R.drawable.topic_image),
            Topic("LITERASI DIGITAL", "Penggunaan Alat Digital dalam Pembelajaran", R.drawable.topic_image2),
            Topic("SOFT SKILLS", "Mengembangkan Pola Pikir Kritis", R.drawable.topic_image3),
            Topic("SOSIAL MORAL", "Membentuk Kebiasaan Bertanggung Jawab pada Anak", R.drawable.topic_image4)
        )
        val recyclerView: RecyclerView = view.findViewById(R.id.discover_rv)
        recyclerView.layoutManager = GridLayoutManager(context, 2)
        recyclerView.adapter = TopicAdapter(topics)

        val recyclerViewBaru: RecyclerView = view.findViewById(R.id.discover_rv2)
        recyclerViewBaru.layoutManager = GridLayoutManager(context, 2)
        recyclerViewBaru.adapter = TopicAdapter(topics)

        getMe()
        setData()
        observeData()
        getProfilePicture()
    }
    private fun getMe(){
        runBlocking{
            val token = userPreference.getToken()
            if (token != null) {
                discoverViewModel.getMe(token)
            }
        }
    }

    private fun setData(){
        discoverViewModel.getMeResponse.observe(viewLifecycleOwner) { response ->
            val id = response.data?.id
            val email = response.data?.email
            val name = response.data?.name
            val nama_instansi = response.data?.asalInstansi
            val tanggal_lahir = response.data?.dateOfBirth
            if (id != null && email != null && name != null) {
                val user = User(id, email, name, glideUrl, nama_instansi, tanggal_lahir)
                sharedViewModel.setUser(user)
            }
        }
    }
    private fun observeData() {
        sharedViewModel.user.observe(viewLifecycleOwner){user->
            val nameDiscover = binding.discoverDiscoverNama
            nameDiscover.text = user.name
        }
    }

    private fun checkToken(){
        runBlocking {
            val token = userPreference.getToken()
            if(token == null){
                val intent = Intent(requireContext(), MainActivity::class.java )
                startActivity(intent)
            }
        }
    }

    private fun getProfilePicture(){
        runBlocking {
            val token = userPreference.getToken()
            val imageUrl = "https://teachtechai.et.r.appspot.com/api/user/profile-picture"
            glideUrl = GlideUrl(
                imageUrl,
                LazyHeaders.Builder()
                    .addHeader("Authorization","Bearer $token")
                    .build()
            )
        }
        Glide.with(this)
            .load(glideUrl)
            .transform(CircleCrop(), CenterCrop())
            .into(binding.discoverProfile)
    }
}